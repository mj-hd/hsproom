package controllers

import (
	"net/http"
	"os"

	"hsproom/config"
//	"hsproom/models"
	"hsproom/templates"
	"hsproom/utils/log"
)

type sourceCreateMember struct {
	*templates.DefaultMember
}

func sourceCreateHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "sourceCreate.tmpl"

	err := tmpl.Render(document, sourceCreateMember{
		DefaultMember: &templates.DefaultMember{
			Title: "ソースコードの作成 - " + config.SiteTitle,
			User: getSessionUser(request),
		},
	})

	if err != nil {
		log.Fatal(os.Stdout, err)
		showError(document, request, "ページの表示に失敗しました。管理人へ問い合わせてください。")
	}

}

type sourceEditMember struct {
	*templates.DefaultMember
}

func sourceEditHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "sourceEdit.tmpl"

	err := tmpl.Render(document, sourceEditMember{
		DefaultMember: &templates.DefaultMember{
			Title: "ソースコードの編集 - " + config.SiteTitle,
			User: getSessionUser(request),
		},
	})

	if err != nil {
		log.Fatal(os.Stdout, err)
		showError(document, request, "ページの表示に失敗しました。管理人へ問い合わせてください。")
	}

}
