package models

import "time"

type Chats struct {
	ID        int
	UserID    int
	CreatedAt time.Time
}
