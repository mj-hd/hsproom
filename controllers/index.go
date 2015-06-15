package controllers

import (
	"net/http"
	"os"

	"../config"
	"../models"
	"../templates"
	"../utils/log"
)

type indexMember struct {
	*templates.DefaultMember
	RecentPrograms *[]models.Program
}

func indexHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template

	tmpl.Layout = "default.tmpl"
	tmpl.Template = "index.tmpl"

	var programs []models.Program

	_, err := models.GetProgramListBy(models.ProgramColCreatedAt, &programs, true, 0, 4)

	if err != nil {
		log.Fatal(os.Stdout, err)

		showError(document, request, "ページの読み込みに失敗しました。")

		return
	}

	err = tmpl.Render(document, indexMember{
		DefaultMember: &templates.DefaultMember{
			Title:  config.SiteTitle,
			UserID: getSessionUser(request),
		},
		RecentPrograms: &programs,
	})
	if err != nil {
		log.Fatal(os.Stdout, err)
		showError(document, request, "ページの表示に失敗しました。")
		return
	}
}
