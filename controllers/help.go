package controllers

import (
	"net/http"
	"os"

	"hsproom/config"
	"hsproom/templates"
	"hsproom/utils/log"
)

func helpHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "help.tmpl"

	err := tmpl.Render(document, &templates.DefaultMember{
		Title: "ヘルプ - " + config.SiteTitle,
		User:  getSessionUser(request),
	})

	if err != nil {
		log.Fatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました。")
	}

}
