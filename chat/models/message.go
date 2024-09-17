package models

import (
	"time"
)

type IncomingMessage struct {
	ID           int       `json:"id"`
	Type         int       `json:"type"`
	Code         string    `json:"code"`
	ChatID       int       `json:"chatId"`
	ChatCode     string    `json:"chatCode"`
	SenderID     int       `json:"senderId"`
	SenderCode   string    `json:"senderCode"`
	ReceiverCode []string  `json:"receiverCode"`
	ReceiverID   int       `json:"receiverId"`
	Seen         bool      `json:"seen"`
	SeenAt       time.Time `json:"seenAt"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Asset        `json:"asset"`
}

type OutgoingMessage struct {
	Code          string    `json:"messageCode"`
	Content       string    `json:"content"`
	Seen          bool      `json:"seen"`
	SeenAt        time.Time `json:"seenAt"`
	Type          int16     `json:"type"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ReceiversCode []string  `json:"code"`
	Asset         struct {
		Path     string `json:"path"`
		MIMEType string `json:"mimeType"`
	} `json:"asset"`
}

type Asset struct {
	ID           int       `json:"id"`
	UserID       int       `json:"userId"`
	Name         string    `json:"name"`
	MimeType     string    `json:"mimeType"`
	Content      string    `json:"content"`
	LastModified time.Time `json:"lastModified"`
	Reason       int16     `json:"reason"`
	Size         int64     `json:"size"`
}

type IncomingMessageResult struct {
	Result chan IncomingMessage
	Err    chan error
}

type OutgoingMessageResult struct {
	Result chan OutgoingMessage
	Err    chan error
}
