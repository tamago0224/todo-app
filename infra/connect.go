package infra

import (
	"database/sql"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

func OpenDB() (*sql.DB, error) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}
	// 環境変数から接続先DBを取得する
	c := mysql.Config{
		DBName:               os.Getenv("TODO_DBNAME"),
		User:                 os.Getenv("TODO_USERNAME"),
		Passwd:               os.Getenv("TODO_PASSWORD"),
		Addr:                 os.Getenv("TODO_HOSTNAME"),
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  jst,
	}
	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}
