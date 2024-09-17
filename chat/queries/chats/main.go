package chats

import "time"

type Message struct {
	ChatID           int       `json:"chatId"`
	ChatCode         string    `json:"chatCode"`
	MessageID        int       `json:"messageId"`
	MessageContent   string    `json:"messageContent"`
	MessageCreatedAt time.Time `json:"messageCreatedAt"`
	Asset            `json:"asset"`
}

type Asset struct {
	ID           int       `json:"id"`
	UserID       int       `json:"userId"`
	Name         string    `json:"name"`
	MimeType     string    `json:"mimeType"`
	Content      string    `json:"content"`
	LastModified time.Time `json:"lastModified"`
	Size         int64     `json:"size"`
	Reason       uint8     `json:"reason"`
}
