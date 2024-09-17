package models

import "time"

type Users struct {
	ID        int
	FirstName string
	LastName  string
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string
	CreatedAt time.Time
}
