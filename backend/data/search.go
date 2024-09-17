package data

import (
	"chatin/auth"
	"chatin/database"
	searchQueries "chatin/queries/search"
	"chatin/utils"
	"database/sql"
	"time"
)

type searchData struct{}

func NewSearchData() *searchData {
	return &searchData{}
}

func (sd *searchData) GetSearchResultRowNumber(searchArgs searchQueries.SearchArguments, searchResultRowNumberChan chan int, errorChan chan error) {
	var (
		statement *sql.Stmt
		err       error
		db        = database.OpenConnection()
		rowNumber = 0
	)
	if searchArgs.ChatID != 0 {
		statement, err = db.Prepare("CALL GET_SEARCH_RESULTS_BY_CHAT_ID_NUMBER(?, ?, ?)")
		if err != nil {
			errorChan <- err
			return
		}

		if statement.QueryRow(&searchArgs.ChatID, &searchArgs.UserID, &searchArgs.Term).Scan(&rowNumber) != nil {
			errorChan <- err
			return
		}
		searchResultRowNumberChan <- rowNumber
		return
	}

	statement, err = db.Prepare("CALL GET_SEARCH_RESULTS_NUMBER(?, ?)")
	if err != nil {
		errorChan <- err
		return
	}

	if statement.QueryRow(&searchArgs.UserID, &searchArgs.Term).Scan(&rowNumber) != nil {
		errorChan <- err
		return
	}

	searchResultRowNumberChan <- rowNumber
}

func (sd *searchData) GetSearchResultsByText(searchArgs searchQueries.SearchArguments, searchChannel searchQueries.SearchResultChannel) {
	var (
		statement          *sql.Stmt
		err                error
		db                 = database.OpenConnection()
		rows               *sql.Rows
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

	statement, err = db.Prepare("CALL GET_SEARCH_MESSAGES_BY_TEXT(?, ?, ?)")
	if err != nil {
		searchChannel.Err <- err
		return
	}

	rows, err = statement.Query(&searchArgs.Term, &searchArgs.MessageID, &searchArgs.UserID)
	if err != nil {
		searchChannel.Err <- err
		return
	}
	var foundItems []searchQueries.Items
	var foundItem searchQueries.Items
	for rows.Next() {
		err := rows.Scan(&chatId, &groupId, &groupName, &messageId, &senderId, &senderFirstName,
			&senderLastName, &recipientId, &recipientFirstName, &recipientLastName, &content, &seen, &createdAt, &assetName, &assetMime)
		if err != nil {
			if err == sql.ErrNoRows {
				searchChannel.FoundItems <- []searchQueries.Items{}
				searchChannel.Err <- nil
				return
			}

			searchChannel.Err <- nil
			return
		}
		foundItem.MessageSeen = utils.GetNullString(seen) == "\x00"
		foundItem.AssetName = new(string)
		*foundItem.AssetName = utils.GetNullString(assetName)
		foundItem.AssetMimeType = new(string)
		*foundItem.AssetMimeType = utils.GetNullString(assetMime)
		foundItem.GroupName = new(string)
		*foundItem.GroupName = utils.GetNullString(groupName)
		foundItem.GroupCode = new(string)
		*foundItem.GroupCode, _ = auth.EncodeUserID(uint64(utils.GetNullInt(groupId)))
		foundItem.ChatCode, _ = auth.EncodeUserID(uint64(chatId))
		foundItem.SenderCode, _ = auth.EncodeUserID(uint64(senderId))
		foundItem.SenderFullName = senderFirstName + " " + senderLastName
		foundItem.RecipientCode, _ = auth.EncodeUserID(uint64(recipientId))
		foundItem.RecipientFullName = recipientFirstName + " " + recipientLastName
		foundItem.MessageContent = content
		foundItem.MessageCreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		foundItems = append(foundItems, foundItem)

	}
	searchChannel.FoundItems <- foundItems

}

func (sd *searchData) GetSearchResultsByChatID(searchArgs searchQueries.SearchArguments, searchChannel searchQueries.SearchResultChannel) {
	var (
		statement          *sql.Stmt
		err                error
		db                 = database.OpenConnection()
		rows               *sql.Rows
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

	statement, err = db.Prepare("CALL GET_SEARCH_MESSAGES_BY_CHAT_ID(?, ?, ?, ?)")
	if err != nil {
		searchChannel.Err <- err
		return
	}

	rows, err = statement.Query(&searchArgs.ChatID, &searchArgs.MessageID, &searchArgs.UserID, &searchArgs.Term)
	if err != nil {
		searchChannel.Err <- err
		return
	}
	var foundItems []searchQueries.Items
	var foundItem searchQueries.Items
	for rows.Next() {
		err := rows.Scan(&groupName, &messageId, &senderId, &senderFirstName,
			&senderLastName, &recipientId, &recipientFirstName, &recipientLastName, &content, &seen, &createdAt, &assetName, &assetMime)
		if err != nil {
			if err == sql.ErrNoRows {
				searchChannel.FoundItems <- []searchQueries.Items{}
				searchChannel.Err <- nil
				return
			}

			searchChannel.Err <- nil
			return
		}
		foundItem.MessageSeen = utils.GetNullString(seen) == "\x00"
		foundItem.AssetName = new(string)
		*foundItem.AssetName = utils.GetNullString(assetName)
		foundItem.AssetMimeType = new(string)
		*foundItem.AssetMimeType = utils.GetNullString(assetMime)
		foundItem.GroupName = new(string)
		*foundItem.GroupName = utils.GetNullString(groupName)
		foundItem.GroupCode = new(string)
		*foundItem.GroupCode, _ = auth.EncodeUserID(uint64(utils.GetNullInt(groupId)))
		foundItem.ChatCode, _ = auth.EncodeUserID(uint64(chatId))
		foundItem.SenderCode, _ = auth.EncodeUserID(uint64(senderId))
		foundItem.SenderFullName = senderFirstName + " " + senderLastName
		foundItem.RecipientCode, _ = auth.EncodeUserID(uint64(recipientId))
		foundItem.RecipientFullName = recipientFirstName + " " + recipientLastName
		foundItem.MessageContent = content
		foundItem.MessageCreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		foundItems = append(foundItems, foundItem)

	}
	searchChannel.FoundItems <- foundItems

}
