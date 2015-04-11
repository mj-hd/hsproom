package models

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"hsproom/config"
	"hsproom/utils/log"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName))
	if err != nil {
		log.Fatal(os.Stdout, err)
		panic(err.Error())
	}
}

func Del() {
	DB.Close()
}
