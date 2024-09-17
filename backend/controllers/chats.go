package controllers

import (
	"chatin/auth"
	"chatin/config/keys"
	messagesQueries "chatin/queries/messages"
	"chatin/server/response"
	"chatin/services"
	"net/http"

	"github.com/gorilla/mux"
)

type chatsControlller struct{}

func NewChatsController() *chatsControlller {
	return &chatsControlller{}
}

var chatsService = services.NewChatsService()

func (c *chatsControlller) Get(res http.ResponseWriter, req *http.Request) {
	req.Context().Value(keys.UserContextKey)
	var params = mux.Vars(req)
	var chatCode = params["chatID"]
	var chatMessages messagesQueries.ChatDetailsData = messagesQueries.ChatDetailsData{
		Rows: make(chan messagesQueries.ChatDetails),
		Err:  make(chan error),
	}
	var id, _ = auth.DecodeUserID(chatCode)
	go chatsService.GetChatDetails(id, chatMessages)
	var err error
	var ChatDetails messagesQueries.ChatDetails
	select {
	case ChatDetails = <-chatMessages.Rows:
		break
	case err = <-chatMessages.Err:
		break
	}
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		http.Error(res, err.Error(), 500)
	}
	response.JSON(res, ChatDetails)

}

func (c *chatsControlller) GetChatList(res http.ResponseWriter, req *http.Request) {
	var id = req.Context().Value(keys.UserContextKey).(float64)
	var ChatListData messagesQueries.ChatListData = messagesQueries.ChatListData{
		Rows: make(chan []messagesQueries.ChatList),
		Err:  make(chan error),
	}
	go chatsService.FindChatList(int(id), ChatListData)
	var chatList []messagesQueries.ChatList
	var err error
	select {
	case chatList = <-ChatListData.Rows:
		break
	case err = <-ChatListData.Err:
		break
	}

	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		http.Error(res, err.Error(), 500)
	}
	response.JSON(res, chatList)
}

func (c *chatsControlller) GetMoreMessages(res http.ResponseWriter, req *http.Request) {
	req.Context().Value(keys.UserContextKey)
	var params = mux.Vars(req)
	var chatCode = params["chatID"]
	var rawMessageId = req.URL.Query().Get("k")
	if rawMessageId == "" {
		res.Header().Set("Content-Type", "application/json")
		http.Error(res, "Key not provided", 400)
	}
	messageId, _ := auth.DecodeUserID(rawMessageId)
	var chatMessages messagesQueries.GetMoreMessages = messagesQueries.GetMoreMessages{
		Rows: make(chan []messagesQueries.Messages),
		Err:  make(chan error),
	}
	var chatId, _ = auth.DecodeUserID(chatCode)
	go chatsService.GetMoreMessages(chatId, messageId, chatMessages)
	var err error
	var messages []messagesQueries.Messages
	select {
	case messages = <-chatMessages.Rows:
		break
	case err = <-chatMessages.Err:
		break
	}
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		http.Error(res, err.Error(), 500)
	}
	response.JSON(res, messages)

}
