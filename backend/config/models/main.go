package models

type Config struct {
	Chatin struct {
		ConnectionString string `json:"connection-string"`
		SecretKey        string `json:"secret-key"`
	} `json:"chatin"`
}
