package controllers

import (
	"net/http"

	"../config"
	"../templates"
)

func aboutHandler(document http.ResponseWriter, request *http.Request) (err error) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "about.tmpl"

	return tmpl.Render(document, &templates.DefaultMember{
		Title:  "このサイトについて - " + config.SiteTitle,
		UserID: getSessionUser(request),
	})
}
