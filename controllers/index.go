package controllers

import (
	"net/http"
	"os"

	"hsproom/config"
	"hsproom/models"
	"hsproom/templates"
	"hsproom/utils"
)

type indexMember struct {
	*templates.DefaultMember
	RecentPrograms *[]models.ProgramInfo
}

func indexHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template

	tmpl.Layout = "default.tmpl"
	tmpl.Template = "index.tmpl"

	var programs []models.ProgramInfo

	_, err := models.GetProgramListBy(models.ProgramColCreated, &programs, true, 0, 4)

	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)

		showError(document, request, "ページの読み込みに失敗しました。")

		return
	}

	err = tmpl.Render(document, indexMember{
		DefaultMember: &templates.DefaultMember{
			Title: config.SiteTitle,
			User:  getSessionUser(request),
		},
		RecentPrograms: &programs,
	})
	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)
		showError(document, request, "ページの表示に失敗しました。")
		return
	}
}
