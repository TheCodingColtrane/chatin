package database

import (
	"chatin/config"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func OpenConnection() *sql.DB {
	db, err := sql.Open("mysql", config.NewConfig().ConnectionString)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db

}
