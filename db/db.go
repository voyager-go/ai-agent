package db

import (
	"ai-agent/shared"
	"database/sql"
	"log"
)

func InitDbConnect() {
	dsn := "root:123456@tcp(127.0.0.1:3307)/test_db"
	var err error
	shared.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
}
