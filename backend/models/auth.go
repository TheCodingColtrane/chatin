package models

type Authentication struct {
	Token  string `json:"token"`
	Expiry int64  `json:"exp"`
}
