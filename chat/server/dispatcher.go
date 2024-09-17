package server

import (
	"log"
	"messenger/auth"
	"messenger/enums"
	"messenger/models"
)

type messageDispatcher struct{}

func NewDispatcher() *messageDispatcher {

	return &messageDispatcher{}
}

func (d *messageDispatcher) Dispatch(action int16,
	incomingMessage models.IncomingMessage,
	outgoingMessageResult models.OutgoingMessageResult) {
	var incomingMessageResult models.IncomingMessageResult
	var outgoingMessage models.OutgoingMessage
	incomingMessageResult.Result = make(chan models.IncomingMessage)
	incomingMessageResult.Err = make(chan error)
	if action == enums.CREATE_MESSAGE {
		go createMessage(incomingMessage, incomingMessageResult)
		select {
		case createdMessage := <-incomingMessageResult.Result:
			outgoingMessage.Content = incomingMessage.Content
			outgoingMessage.CreatedAt = createdMessage.CreatedAt
			outgoingMessage.ReceiversCode = incomingMessage.ReceiverCode
			outgoingMessage.Asset.Path = incomingMessage.Asset.Name
			outgoingMessage.Asset.MIMEType = incomingMessage.Asset.MimeType
			outgoingMessage.Code = createdMessage.Code
			outgoingMessage.Type = 0
			outgoingMessage.Seen = false
			outgoingMessageResult.Result <- outgoingMessage
			break
		case err := <-incomingMessageResult.Err:
			outgoingMessageResult.Err <- err
			break
		}
	} else if action == enums.UPDATE_MESSAGE {
		go updateMessage(incomingMessage, incomingMessageResult)
		var chatAction int16
		if incomingMessage.Type == 5 {
			chatAction = 1
		} else {
			chatAction = 2
		}
		select {

		case incomingMessage = <-incomingMessageResult.Result:
			outgoingMessage.Content = incomingMessage.Content
			outgoingMessage.Code = incomingMessage.Code
			outgoingMessage.UpdatedAt = incomingMessage.UpdatedAt
			outgoingMessage.ReceiversCode = incomingMessage.ReceiverCode
			outgoingMessage.Seen = true
			outgoingMessage.SeenAt = incomingMessage.SeenAt
			outgoingMessage.Type = chatAction
			outgoingMessageResult.Result <- outgoingMessage
			break
		case err := <-incomingMessageResult.Err:
			outgoingMessageResult.Err <- err
			break
		}

	}

}
func createMessage(incomingMessage models.IncomingMessage, createdMessage models.IncomingMessageResult) {
	var size int64
	var outgoingMessage models.OutgoingMessage
	if incomingMessage.Type > 1 {

		createdFile, err := CreateAsset(incomingMessage)
		if err != nil {
			log.Printf("message service error: %v", err)
			return
		}
		outgoingMessage.Asset.MIMEType = createdFile.MimeType
		outgoingMessage.Asset.Path = createdFile.FileName[1:]
		size = createdFile.Size

	}

	id, _ := auth.DecodeID(incomingMessage.ChatCode)
	incomingMessage.ChatID = id

	go messageService.Create(models.IncomingMessage{
		SenderID: incomingMessage.SenderID,
		Content:  incomingMessage.Content,
		ChatCode: incomingMessage.ChatCode,
		ChatID:   incomingMessage.ChatID,
		Asset: models.Asset{
			Name: outgoingMessage.Asset.Path, MimeType: outgoingMessage.Asset.MIMEType, Size: size,
			UserID: incomingMessage.SenderID,
		}}, createdMessage)

	select {
	case message := <-createdMessage.Result:
		createdMessage.Result <- message
	case err := <-createdMessage.Err:
		createdMessage.Err <- err

	}
}

func updateMessage(incomingMessage models.IncomingMessage, updatedMessage models.IncomingMessageResult) {
	//var ids = strings.Split(incomingMessage.Code, "-")
	incomingMessage.ID, _ = auth.DecodeID(incomingMessage.Code)
	go messageService.UpdateMessage(incomingMessage, updatedMessage)
	select {
	case message := <-updatedMessage.Result:
		updatedMessage.Result <- message
	case err := <-updatedMessage.Err:
		updatedMessage.Err <- err

	}
}
