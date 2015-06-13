package controllers

import (
	"net/http"
	"os"
	"strconv"

	"../config"
	"../models"
	"../templates"
	"../utils/log"
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
	Program *models.Program
}

func sourceEditHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "sourceEdit.tmpl"

	rawProgramId := request.URL.Query().Get("p")
	programId, err := strconv.Atoi(rawProgramId)

	if err != nil {
		log.Debug(os.Stdout, err)

		showError(document, request, "プログラムが見つかりません。")

		return
	}

	user := getSessionUser(request)

	program := models.NewProgram()
	err = program.Load(programId)

	if err != nil {
		log.Debug(os.Stdout, err)

		showError(document, request, "プログラムの読み込みに失敗しました。")

		return
	}

	if program.User != user {
		log.DebugStr(os.Stdout, "権限のない編集画面へのアクセス")

		showError(document, request, "プログラムの編集権限がありません。")

		return
	}

	err = tmpl.Render(document, sourceEditMember{
		DefaultMember: &templates.DefaultMember{
			Title: "ソースコードの編集 - " + config.SiteTitle,
			User: getSessionUser(request),
		},
		Program: program,
	})

	if err != nil {
		log.Fatal(os.Stdout, err)
		showError(document, request, "ページの表示に失敗しました。管理人へ問い合わせてください。")
	}

}
