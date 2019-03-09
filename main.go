package main

import (
	"github.com/jinzhu/gorm"
	"github.com/simon987/ws_bucket/api"
	"os"
)

func main() {

	db, err := gorm.Open(getDialect(), getConnStr())
	if err != nil {
		panic(err)
	}

	a := api.New(db)
	a.Run()
}

func getConnStr() string {
	connStr := os.Getenv("WS_BUCKET_CONNSTR")
	if connStr == "" {
		return "host=localhost user=ws_bucket dbname=ws_bucket password=ws_bucket sslmode=disable"
	} else {
		return connStr
	}
}

func getDialect() string {
	connStr := os.Getenv("WS_BUCKET_DIALECT")
	if connStr == "" {
		return "postgres"
	} else {
		return connStr
	}
}
