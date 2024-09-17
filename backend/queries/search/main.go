package searchQueries

import "time"

type FoundItemsData struct {
	FoundItems int `json:"foundItems"`
}

type SearchResultData struct {
	Items []Items `json:"items"`
}

type SearchResultChannel struct {
	FoundItems chan []Items
	Err        chan error
}

type SearchArguments struct {
	UserID      float64
	ChatID      int
	ChatCode    string `json:"chatCode"`
	MessageID   int
	MessageCode string `json:"messageCode"`
	Term        string `json:"term"`
}

type Items struct {
	UserCode          string    `json:"userCode"`
	ChatCode          string    `json:"chatCode"`
	GroupCode         *string   `json:"groupCode"`
	GroupName         *string   `json:"groupName"`
	SenderCode        string    `json:"senderCode"`
	SenderFullName    string    `json:"senderFullName"`
	RecipientCode     string    `json:"recipientCode"`
	RecipientFullName string    `json:"recipientFullName"`
	MessageContent    string    `json:"messageContent"`
	MessageSeen       bool      `json:"messageSeen"`
	MessageCreatedAt  time.Time `json:"messageCreatedAt"`
	AssetName         *string   `json:"assetName"`
	AssetMimeType     *string   `json:"assetMimeType"`
}
