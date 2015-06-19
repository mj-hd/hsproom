package controllers

import (
	"net/http"

	"../config"
	"../models"
	"../templates"
)

type indexMember struct {
	*templates.DefaultMember
	RecentPrograms *[]models.Program
}

func indexHandler(document http.ResponseWriter, request *http.Request) (err error) {

	var tmpl templates.Template

	tmpl.Layout = "default.tmpl"
	tmpl.Template = "index.tmpl"

	var programs []models.Program

	_, err = models.GetProgramListBy(models.ProgramColCreatedAt, &programs, true, 0, 4)

	if err != nil {
		return err
	}

	return tmpl.Render(document, indexMember{
		DefaultMember: &templates.DefaultMember{
			Title:  config.SiteTitle,
			UserID: getSessionUser(request),
		},
		RecentPrograms: &programs,
	})
}
