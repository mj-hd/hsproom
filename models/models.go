package models

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"../config"
	"../utils/log"
)

var DB *gorm.DB

func Init() {
	var err error

	DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName))
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}

	if config.LogLevel == log.Level_Debug {
		DB.LogMode(true)
	}

	// 以下にテーブルの初期化処理を書く
	// 新しいテーブルを作成するときは、以下に初期化処理を追記する
	initPrograms()
	initUsers()
	initGoods()
	initComments()
	initNotifications()
	initRelevances()
}

func Del() {
	delRelevances()
	DB.Close()
}
