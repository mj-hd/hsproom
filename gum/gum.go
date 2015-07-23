package gum

import (
	"net/http"

	"../bot"
	"../config"
	"../controllers"
	"../models"
	"../plugins"
	"../templates"
	"../utils/log"

	"github.com/gorilla/context"
)

func init() {

	log.LogFile = config.LogFile
	log.DisplayLog = config.DisplayLog
	log.LogLevel = config.LogLevel

}

func Del() {
	models.Del()
	controllers.Del()
	templates.Del()
	plugins.Del()
	bot.Del()
}

func Start() {
	bot.Init()
	controllers.Init()
	models.Init()
	plugins.Init()
	templates.Init()

	for route := range controllers.Router.Iterator() {
		http.HandleFunc(route.Path, route.Function)
		log.DebugStr(route.Path + "に関数を割当")
	}

	http.Handle("/"+config.StaticPath, http.StripPrefix("/"+config.StaticPath, http.FileServer(http.Dir(config.StaticPath))))
	log.DebugStr("/" + config.StaticPath + "に静的コンテンツを割当")

	log.InfoStr("ポート" + config.ServerPort + "でサーバを開始...")
	http.ListenAndServe(":"+config.ServerPort, context.ClearHandler(http.DefaultServeMux))
}
