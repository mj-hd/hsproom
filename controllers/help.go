package controllers

import (
	"net/http"
	"os"
	"io/ioutil"
	"html/template"

	"hsproom/config"
	"hsproom/templates"
	"hsproom/utils/log"

	"github.com/russross/blackfriday"
)

type helpMember struct {
	*templates.DefaultMember
	HelpContent template.HTML
}

func helpHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "help.tmpl"

	mdFile := config.TemplatesPath + "help/main.md"

	mdRaw, err := ioutil.ReadFile(mdFile)
	if err != nil {
		log.Fatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	markdown := blackfriday.MarkdownCommon(mdRaw)

	err = tmpl.Render(document, helpMember{
		DefaultMember: &templates.DefaultMember{
			Title: "ヘルプ - " + config.SiteTitle,
			User:  getSessionUser(request),
		},
		HelpContent: template.HTML(markdown),
	})

	if err != nil {
		log.Fatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました。")
	}

}
