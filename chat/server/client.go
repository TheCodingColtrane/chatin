package server

import (
	"fmt"
	"log"
	"messenger/auth"
	"messenger/enums"
	"messenger/models"
	"messenger/service"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1000000
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan models.OutgoingMessage
	UserID   int
	UserCode string
}

var messageService = service.NewMessageService()
var dispatcher = NewDispatcher()

func (c *Client) ReadPump(done chan struct{}, once *sync.Once) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		once.Do(func() {
			close(done)
		})
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetWriteDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		var incomingMessageResult models.IncomingMessageResult
		var incomingMessage *models.IncomingMessage
		var outgoingMessage models.OutgoingMessage
		var outgoingMessageResult models.OutgoingMessageResult

		err := c.conn.ReadJSON(&incomingMessage)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		if incomingMessage.ReceiverCode[0] == "" {
			continue
		}

		incomingMessageResult.Result = make(chan models.IncomingMessage)
		incomingMessageResult.Err = make(chan error)
		outgoingMessageResult.Result = make(chan models.OutgoingMessage)
		outgoingMessageResult.Err = make(chan error)
		defer close(incomingMessageResult.Result)
		defer close(incomingMessageResult.Err)
		defer close(outgoingMessageResult.Result)
		defer close(outgoingMessageResult.Err)

		incomingMessage.SenderID = c.UserID
		if incomingMessage.Code == "" {
			go dispatcher.Dispatch(enums.CREATE_MESSAGE, *incomingMessage, outgoingMessageResult)
		} else {
			go dispatcher.Dispatch(enums.UPDATE_MESSAGE, *incomingMessage, outgoingMessageResult)
		}

		select {
		case outgoingMessage = <-outgoingMessageResult.Result:
			c.hub.broadcast <- outgoingMessage
			continue
		case <-outgoingMessageResult.Err:
			c.hub.broadcast <- models.OutgoingMessage{}
			continue
		}

	}
}

func (c *Client) WritePump(done chan struct{}, once *sync.Once) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		once.Do(func() {
			close(done)
		})
	}()
	for {
		select {
		case <-done:
			return

		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("Authorization")
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	claims, err := auth.VerifyToken(token.Value)
	if err != nil || claims == nil {
		fmt.Print(err)
		conn.Close()
		return
	}
	id, _ := strconv.Atoi(claims.UserID)
	code, _ := auth.EncodeID(uint64(id))

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan models.OutgoingMessage, 256),
		UserID:   id,
		UserCode: code,
	}
	client.hub.register <- client

	done := make(chan struct{})
	var once sync.Once
	go client.WritePump(done, &once)
	go client.ReadPump(done, &once)
}
