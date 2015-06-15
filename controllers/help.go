package controllers

import (
	"html/template"
	"net/http"
	"os"

	"../config"
	"../templates"
	"../utils/log"
)

type helpMember struct {
	*templates.DefaultMember
	HelpContent template.HTML
}

func helpHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "help.tmpl"

	err := tmpl.Render(document, helpMember{
		DefaultMember: &templates.DefaultMember{
			Title:  "ヘルプ - " + config.SiteTitle,
			UserID: getSessionUser(request),
		},
	})

	if err != nil {
		log.Fatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました。")
	}
}
