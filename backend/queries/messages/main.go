package messagesQueries

import "time"

type Messages struct {
	Code       string     `json:"code"`
	SenderCode string     `json:"senderCode"`
	Content    string     `json:"content"`
	Seen       bool       `json:"seen"`
	SeenAt     time.Time  `json:"seenAt"`
	CreatedAt  string     `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
	Assets     `json:"asset"`
}

type Assets struct {
	Name     string `json:"name"`
	MimeType string `json:"mimeType"`
	Reason   uint8  `json:"reason"`
}

type Participants struct {
	Code       string `json:"code"`
	FullName   string `json:"fullName"`
	UserName   string `json:"username"`
	ProfileIMG string `json:"profileImg"`
}

type ChatList struct {
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

type ChatDetails struct {
	Messages     []Messages     `json:"messages"`
	Participants []Participants `json:"participants"`
}

type ChatMessages struct {
	Rows chan []Messages
	Err  chan error
}

type ChatParticipants struct {
	Rows chan []Participants
	Err  chan error
}

type ChatDetailsData struct {
	Rows chan ChatDetails
	Err  chan error
}

type ChatListData struct {
	Rows chan []ChatList
	Err  chan error
}

type GetMoreMessages struct {
	Rows chan []Messages
	Err  chan error
}
