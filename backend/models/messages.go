package models

import "time"

type Messages struct {
	ID        int
	UserID    int
	ChatID    int
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
