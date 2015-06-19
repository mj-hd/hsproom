package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"../config"
	"../models"
	"../templates"
	"../utils/log"
)

type sourceCreateMember struct {
	*templates.DefaultMember
	ThumbnailLimitSize   int
	StartaxLimitSize     int
	AttachmentsLimitSize int
}

func sourceCreateHandler(document http.ResponseWriter, request *http.Request) (err error) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "sourceCreate.tmpl"

	return tmpl.Render(document, sourceCreateMember{
		DefaultMember: &templates.DefaultMember{
			Title:  "ソースコードの作成 - " + config.SiteTitle,
			UserID: getSessionUser(request),
		},
		ThumbnailLimitSize:   config.ThumbnailLimitSize,
		StartaxLimitSize:     config.StartaxLimitSize,
		AttachmentsLimitSize: config.AttachmentsLimitSize,
	})
}

type sourceEditMember struct {
	*templates.DefaultMember
	Program              *models.Program
	ThumbnailLimitSize   int
	StartaxLimitSize     int
	AttachmentsLimitSize int
}

func sourceEditHandler(document http.ResponseWriter, request *http.Request) (err error) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "sourceEdit.tmpl"

	rawProgramId := request.URL.Query().Get("p")
	programId, err := strconv.Atoi(rawProgramId)

	if err != nil {
		log.Debug(err)

		showError(document, request, "プログラムが見つかりません。")

		return nil
	}

	user := getSessionUser(request)

	program := models.NewProgram()
	err = program.Load(programId)

	if err != nil {
		return errors.New("プログラムの読み込みに失敗: \r\n" + err.Error())
	}

	if program.UserID != user {
		log.DebugStr("権限のない編集画面へのアクセス")

		showError(document, request, "プログラムの編集権限がありません。")

		return nil
	}

	err = program.LoadThumbnail()
	if err != nil {
		log.DebugStr("サムネイル画像の読み込みに失敗しました。")

		showError(document, request, "サムネイル画像の読み込みに失敗しました。")

		return nil
	}

	err = program.LoadAttachments()
	if err != nil {
		log.DebugStr("添付ファイルの読み込みに失敗しました。")

		showError(document, request, "添付ファイルの読み込みに失敗しました。")

		return nil
	}

	return tmpl.Render(document, sourceEditMember{
		DefaultMember: &templates.DefaultMember{
			Title:  "ソースコードの編集 - " + config.SiteTitle,
			UserID: getSessionUser(request),
		},
		Program:              program,
		ThumbnailLimitSize:   config.ThumbnailLimitSize,
		StartaxLimitSize:     config.StartaxLimitSize,
		AttachmentsLimitSize: config.AttachmentsLimitSize,
	})
}
