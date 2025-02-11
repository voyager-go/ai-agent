package db

import (
	"ai-agent/shared"
	"database/sql"
	"log"
	"os"
)

func InitDbConnect() {
	dsn := os.Getenv("DB_DSN")
	var err error
	shared.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
}
