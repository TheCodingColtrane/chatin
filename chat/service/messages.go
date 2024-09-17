package service

import (
	"database/sql"
	"messenger/auth"
	"messenger/database"
	"messenger/enums"
	"messenger/models"
	"time"
)

type messageService struct{}

func NewMessageService() *messageService {
	return &messageService{}
}
func (s *messageService) Create(message models.IncomingMessage, createdMessage models.IncomingMessageResult) {
	var db = database.OpenConnection()
	var asset models.IncomingMessage
	var err error
	if message.Asset.Name != "" {
		go s.CreateAsset(message, createdMessage)
		select {
		case asset = <-createdMessage.Result:
			break
		case <-createdMessage.Err:
			return
		}
	}
	var statement *sql.Stmt
	if asset.Asset.ID > 0 {
		statement, err = db.Prepare("INSERT INTO messages (chat_id, user_id, asset_id, content) VALUES (?,?,?, ?)")
	} else {
		statement, err = db.Prepare("INSERT INTO messages (chat_id, user_id, content) VALUES (?,?,?)")

	}
	if err != nil {
		createdMessage.Err <- err
		return
	}
	var result sql.Result
	if asset.Asset.ID > 0 {
		result, err = statement.Exec(&message.ChatID, message.SenderID, &asset.Asset.ID, &message.Content)

	} else {
		result, err = statement.Exec(&message.ChatID, message.SenderID, &message.Content)
	}
	if err != nil {
		createdMessage.Err <- err
		return
	}

	id, _ := result.LastInsertId()
	if id > 0 {
		message.CreatedAt = time.Now()
		message.ID = int(id)
		message.Code, _ = auth.EncodeID(uint64(id))
		createdMessage.Result <- message
		return
	}

}

func (s *messageService) CreateAsset(asset models.IncomingMessage, createdMessage models.IncomingMessageResult) {
	var db = database.OpenConnection()
	statement, err := db.Prepare("INSERT INTO assets (user_id, name, mime_type, size, reason) VALUES (?,?,?, ?, ?)")
	if err != nil {
		createdMessage.Err <- err
		return
	}
	var mime = 0
	switch asset.MimeType {
	case "audio/webm":
		mime = enums.AUDIO
	case "image/jpeg":
		mime = enums.IMAGE
	case "image/png":
		mime = enums.IMAGE
	case "video/mp4":
		mime = enums.VIDEO
	}

	result, err := statement.Exec(&asset.UserID, &asset.Name, &mime, &asset.Size, &asset.Reason)
	if err != nil {
		createdMessage.Err <- err
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		createdMessage.Err <- err
		return
	}

	asset.ID = int(id)
	createdMessage.Result <- asset
	select {
	case message := <-createdMessage.Result:
		createdMessage.Result <- message
	case err := <-createdMessage.Err:
		createdMessage.Err <- err

	}

}

func (s *messageService) UpdateMessage(message models.IncomingMessage, messageResult models.IncomingMessageResult) {

	var db = database.OpenConnection()
	var sqlStatement = ""
	if message.Content != "" {
		sqlStatement = "UPDATE messages SET content = ?,  updated_at = ? WHERE id = ?"
	} else {
		sqlStatement = "UPDATE messages SET seen = ?, seen_at = ? WHERE id = ?"
	}

	statement, err := db.Prepare(sqlStatement)
	if err != nil {
		messageResult.Err <- err
		return
	}

	var result sql.Result
	if message.Content != "" {
		message.UpdatedAt = time.Now()
		result, err = statement.Exec(&message.Content, &message.UpdatedAt, &message.ID)

	} else {
		message.Seen = true
		message.SeenAt = time.Now()
		result, err = statement.Exec(&message.Seen, &message.SeenAt, &message.ID)

	}
	if err != nil {
		messageResult.Err <- err
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		messageResult.Err <- err
		return
	}

	if rowsAffected > 0 {
		messageResult.Result <- message
		return
	}

}
