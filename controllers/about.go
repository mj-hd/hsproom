package controllers

import (
	"net/http"
	"os"

	"hsproom/config"
	"hsproom/templates"
	"hsproom/utils"
)

func aboutHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "about.tmpl"

	err := tmpl.Render(document, &templates.DefaultMember{
		Title: "このサイトについて - " + config.SiteTitle,
		User:  getSessionUser(request),
	})

	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました。")
	}

}
