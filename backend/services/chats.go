package services

import (
	"chatin/auth"
	"chatin/database"
	messagesQueries "chatin/queries/messages"
	"chatin/utils"
	"database/sql"
	"fmt"
	"time"
)

type chatsService struct{}

func NewChatsService() *chatsService {
	return &chatsService{}
}

func (s *chatsService) GetChatDetails(id int, chatMessages messagesQueries.ChatDetailsData) {
	defer close(chatMessages.Rows)
	defer close(chatMessages.Err)
	var db = database.OpenConnection()
	statement, err := db.Prepare("CALL GET_CHAT_DETAILS(?)")
	if err != nil {
		chatMessages.Err <- err
		return
	}
	rows, err := statement.Query(&id)
	if err == sql.ErrNoRows {
		chatMessages.Err <- nil
		return
	}
	var message messagesQueries.Messages
	var messages = make([]messagesQueries.Messages, 0)
	var chatParticipant messagesQueries.Participants
	var chatParticipants = make([]messagesQueries.Participants, 0)
	var userId = 0
	var mime sql.NullInt16
	var reason sql.NullInt16
	var name sql.NullString
	var seen sql.NullString
	var seenAt sql.NullString
	var updatedAt sql.NullString
	var messageId int
	for rows.Next() {
		err := rows.Scan(&messageId, &message.Content, &message.CreatedAt, &updatedAt, &seen, &seenAt, &userId, &mime, &reason, &name)
		if err != nil {
			fmt.Print(err)
			chatMessages.Err <- err
			return
		}
		if seen.Valid {
			if seen.String != "\x00" {
				message.Seen = true
			}
			if seenAt.Valid && seenAt.String != "" {
				message.SeenAt, _ = time.Parse("2006-01-02 15:04:05", seenAt.String)
			}

		}
		if mime.Valid {
			message.Assets.MimeType = utils.GetMIMETypeString(uint8(mime.Int16))
			message.Reason = uint8(reason.Int16)
			message.Assets.Name = name.String[1:]

		}

		if updatedAt.Valid {
			var now = time.Now()
			message.UpdatedAt = &now
			*message.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt.String)
		}

		message.Code, _ = auth.EncodeUserID(uint64(messageId))
		message.SenderCode, _ = auth.EncodeUserID(uint64(userId))
		messages = append(messages, message)

		// chatMessages.Rows <- contacts

	}

	if rows.NextResultSet() {
		var (
			firstName     = ""
			lastName      = ""
			participantId = 0
		)
		for rows.Next() {
			err := rows.Scan(&participantId, &firstName, &lastName, &chatParticipant.UserName, &chatParticipant.ProfileIMG)
			if err != nil {
				chatMessages.Err <- err
				return
			}
			chatParticipant.Code, _ = auth.EncodeUserID(uint64(participantId))
			chatParticipant.FullName = firstName + " " + lastName
			chatParticipants = append(chatParticipants, chatParticipant)
		}

	}
	chatMessages.Rows <- messagesQueries.ChatDetails{
		Messages:     messages,
		Participants: chatParticipants,
	}

}

func (s *chatsService) GetMoreMessages(chatId int, messageId int, chatMessages messagesQueries.GetMoreMessages) {
	defer close(chatMessages.Rows)
	defer close(chatMessages.Err)
	var db = database.OpenConnection()
	statement, err := db.Prepare("CALL GET_MORE_MESSAGES(?, ?)")
	if err != nil {
		chatMessages.Err <- err
		return
	}
	rows, err := statement.Query(&chatId, &messageId)
	if err == sql.ErrNoRows {
		chatMessages.Err <- nil
		chatMessages.Rows <- []messagesQueries.Messages{}
		return
	}
	var message messagesQueries.Messages
	var messages = make([]messagesQueries.Messages, 0)
	var userId = 0
	var mime sql.NullInt16
	var reason sql.NullInt16
	var name sql.NullString
	var seen sql.NullString
	var seenAt sql.NullString
	var updatedAt sql.NullString
	var foundMessageId int
	for rows.Next() {
		err := rows.Scan(&foundMessageId, &message.Content, &message.CreatedAt, &updatedAt, &seen, &seenAt, &userId, &mime, &reason, &name)
		if err != nil {
			fmt.Print(err)
			chatMessages.Err <- err
			return
		}
		if seen.Valid {
			if seen.String != "\x00" {
				message.Seen = true
			}
			if seenAt.Valid && seenAt.String != "" {
				message.SeenAt, _ = time.Parse("2006-01-02 15:04:05", seenAt.String)
			}

		}
		if mime.Valid {
			message.Assets.MimeType = utils.GetMIMETypeString(uint8(mime.Int16))
			message.Reason = uint8(reason.Int16)
			message.Assets.Name = name.String[1:]

		}

		if updatedAt.Valid {
			*message.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt.String)
		}

		message.Code, _ = auth.EncodeUserID(uint64(foundMessageId))
		message.SenderCode, _ = auth.EncodeUserID(uint64(userId))
		messages = append(messages, message)

		// chatMessages.Rows <- contacts

	}

	chatMessages.Rows <- messages

}

func (s *chatsService) FindParticipants(id int, result messagesQueries.ChatParticipants) {
	var db = database.OpenConnection()
	statement, err := db.Prepare("SELECT u.id, u.first_name, u.last_name, u.username, a.name FROM chat_participants c " +
		"JOIN users u ON c.user_id = u.id LEFT JOIN assets a ON c.user_id = a.user_id WHERE c.chat_id = ? and a.reason = 5")
	if err != nil {
		result.Err <- err
		return
	}
	rows, err := statement.Query(&id)

	if err != nil {
		result.Err <- err
		return
	}
	var chatParticipant messagesQueries.Participants
	var chatParticipants = make([]messagesQueries.Participants, 0)
	var (
		firstName     = ""
		lastName      = ""
		participantId = 0
	)
	for rows.Next() {
		err := rows.Scan(&participantId, &firstName, &lastName, &chatParticipant.UserName, &chatParticipant.ProfileIMG)
		if err != nil {
			result.Err <- err
			return
		}
		chatParticipant.Code, _ = auth.EncodeUserID(uint64(participantId))
		chatParticipant.FullName = firstName + " " + lastName
		chatParticipants = append(chatParticipants, chatParticipant)

	}

	result.Rows <- chatParticipants

}

func (s *chatsService) FindChatList(id int, chatListData messagesQueries.ChatListData) {
	var db = database.OpenConnection()
	var statement, err = db.Prepare("select c.id \"chat_id\", cg.id \"group_id\", cg.name \"group_name\", m.id \"message_id\", u.id \"sender_id\"," +
		"u.first_name \"sender_first_name\", u.last_name \"sender_last_name\", u2.id \"recipient_id\", u2.first_name \"recipient_first_name\", " +
		"u2.last_name \"recipient_last_name\", m.content, m.seen, m.created_at, a.name, a.mime_type from chats c left join chat_groups cg on " +
		"cg.chat_id = c.id left join assets a on cg.asset_id = a.id join chat_participants cp on cp.chat_id = c.id join messages m on " +
		"c.id = m.chat_id join users u on m.user_id = u.id join users u2 on c.user_id = u2.id where cp.user_id = ? and m.id = " +
		"(select max(m2.id) from messages m2 where m2.chat_id = c.id )")
	if err != nil {
		chatListData.Err <- err
		return
	}

	rows, err := statement.Query(&id)
	if err != nil {
		chatListData.Err <- err
		return
	}

	var (
		chatId             = 0
		senderId           = 0
		messageId          = 0
		senderFirstName    = ""
		senderLastName     = ""
		recipientId        = 0
		recipientFirstName = ""
		recipientLastName  = ""
		content            = ""
		createdAt          string
		seen               sql.NullString
		assetName          sql.NullString
		assetMime          sql.NullString
		groupName          sql.NullString
		groupId            sql.NullInt32
	)

	var chatList messagesQueries.ChatList
	var chats = make([]messagesQueries.ChatList, 0)
	for rows.Next() {
		err := rows.Scan(&chatId, &groupId, &groupName, &messageId, &senderId, &senderFirstName,
			&senderLastName, &recipientId, &recipientFirstName, &recipientLastName, &content, &seen, &createdAt, &assetName, &assetMime)
		if err != nil {
			if err == sql.ErrNoRows {
				chatListData.Rows <- []messagesQueries.ChatList{}
				chatListData.Err <- nil
				return
			}

			chatListData.Err <- err
			return
		}
		fmt.Print(utils.GetNullString(assetName))
		chatList.MessageSeen = utils.GetNullString(seen) == "\x00"
		chatList.AssetName = new(string)
		*chatList.AssetName = utils.GetNullString(assetName)
		chatList.AssetMimeType = new(string)
		*chatList.AssetMimeType = utils.GetNullString(assetMime)
		chatList.GroupName = new(string)
		*chatList.GroupName = utils.GetNullString(groupName)
		chatList.GroupCode = new(string)
		*chatList.GroupCode, _ = auth.EncodeUserID(uint64(utils.GetNullInt(groupId)))
		chatList.ChatCode, _ = auth.EncodeUserID(uint64(chatId))
		chatList.SenderCode, _ = auth.EncodeUserID(uint64(senderId))
		chatList.SenderFullName = senderFirstName + " " + senderLastName
		chatList.RecipientCode, _ = auth.EncodeUserID(uint64(recipientId))
		chatList.RecipientFullName = recipientFirstName + " " + recipientLastName
		chatList.MessageContent = content
		chatList.MessageCreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		chats = append(chats, chatList)

	}

	chatListData.Rows <- chats

}
