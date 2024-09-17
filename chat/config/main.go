package config

import (
	"log"
	"os"

	"github.com/dotenv-org/godotenvvault"
)

type Configs struct {
	ConnectionString string
	SecretKey        string
	SquidsAlpha      string
}

func NewConfig() *Configs {
	var err = godotenvvault.Load()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	var connectionString = os.Getenv("DBUSER") + ":" + os.Getenv("DBPWD") + "@tcp(" + os.Getenv("DBHOST") + ":" + os.Getenv("DBPORT") + ")/" + os.Getenv("DBNAME")
	return &Configs{ConnectionString: connectionString, SecretKey: os.Getenv("SSK"), SquidsAlpha: os.Getenv("SQIDSAPLHA")}

}
