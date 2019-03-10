package test

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/simon987/ws_bucket/api"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {

	//db, err := gorm.Open("postgres", "host=localhost user=ws_bucket dbname=ws_bucket password=ws_bucket sslmode=disable")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}

	a := api.New(db)
	go a.Run()

	time.Sleep(time.Millisecond * 100)

	m.Run()
	_ = os.Remove("test.db")
}
