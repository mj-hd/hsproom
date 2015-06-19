package controllers

import (
	"html/template"
	"net/http"

	"../config"
	"../templates"
)

type helpMember struct {
	*templates.DefaultMember
	HelpContent template.HTML
}

func helpHandler(document http.ResponseWriter, request *http.Request) (err error) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "help.tmpl"

	return tmpl.Render(document, helpMember{
		DefaultMember: &templates.DefaultMember{
			Title:  "ヘルプ - " + config.SiteTitle,
			UserID: getSessionUser(request),
		},
	})
}
