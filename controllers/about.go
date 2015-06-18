package controllers

import (
	"net/http"

	"../config"
	"../templates"
	"../utils/log"
)

func aboutHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "about.tmpl"

	err := tmpl.Render(document, &templates.DefaultMember{
		Title:  "このサイトについて - " + config.SiteTitle,
		UserID: getSessionUser(request),
	})

	if err != nil {
		log.Fatal(err)

		showError(document, request, "エラーが発生しました。")
	}

}
