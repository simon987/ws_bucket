package main

import (
	"github.com/jinzhu/gorm"
	"github.com/simon987/ws_bucket/api"
)

func main() {

	db, err := gorm.Open("postgres", "host=localhost user=ws_bucket dbname=ws_bucket password=ws_bucket sslmode=disable")
	if err != nil {
		panic(err)
	}

	a := api.New(db)
	a.Run()
}
