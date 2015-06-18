package controllers

import (
	"net/http"
	"strconv"

	"../config"
	"../models"
	"../templates"
	"../utils/log"
)

type userViewMember struct {
	*templates.DefaultMember
	UserInfo     *models.User
	UserPrograms *[]models.Program
}

func userViewHandler(document http.ResponseWriter, request *http.Request) {

	var err error

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "userView.tmpl"

	rawUid := request.URL.Query().Get("u")
	uid, err := strconv.Atoi(rawUid)

	if err != nil {
		log.Debug(err)

		showError(document, request, "ユーザが見つかりませんでした。")
		return
	}

	var user models.User
	err = user.Load(uid)

	if err != nil {
		log.Debug(err)

		showError(document, request, "ユーザが見つかりませんでした。")
		return
	}

	var programs []models.Program

	_, err = models.GetProgramListByUser(models.ProgramColCreatedAt, &programs, user.ID, true, 0, 4)

	if err != nil {
		log.Fatal(err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	err = tmpl.Render(document, userViewMember{
		DefaultMember: &templates.DefaultMember{
			Title:  user.Name + " のプロフィール - " + config.SiteTitle,
			UserID: getSessionUser(request),
		},
		UserInfo:     &user,
		UserPrograms: &programs,
	})
	if err != nil {
		log.Fatal(err)

		showError(document, request, "ページの表示に失敗しました。管理人へ問い合わせてください。")
	}
}

func userListHandler(document http.ResponseWriter, request *http.Request) {

}

func userLoginHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "userLogin.tmpl"

	err := tmpl.Render(document, templates.DefaultMember{
		Title:  "ログイン",
		UserID: getSessionUser(request),
	})

	if err != nil {
		log.Fatal(err)

		showError(document, request, "ページの表示に失敗しました。管理人へ問い合わせてください。")
	}
}

func userLogoutHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "userLogout.tmpl"

	removeSession(document, request)

	err := tmpl.Render(document, templates.DefaultMember{
		Title:  "ログアウト中です...",
		UserID: 0,
	})
	if err != nil {
		log.Fatal(err)

		showError(document, request, "ログアウト中にエラーが発生しました。")
	}
}

func userEditHandler(document http.ResponseWriter, request *http.Request) {

}

type userProgramsMember struct {
	*templates.DefaultMember
	Programs     []models.Program
	ProgramCount int
	CurPage      int
	MaxPage      int
	Sort         string
	UserName     string
	UserId       int
}

func userProgramsHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "userPrograms.tmpl"

	userId, err := strconv.Atoi(request.URL.Query().Get("u"))

	if err != nil {
		log.Debug(err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	sort := request.URL.Query().Get("s")

	var sortKey models.ProgramColumn
	switch sort {
	case "c":
		sortKey = models.ProgramColCreatedAt
	case "g":
		sortKey = models.ProgramColGood
	case "n":
		sortKey = models.ProgramColTitle
	default:
		sortKey = models.ProgramColCreatedAt
	}

	page, err := strconv.Atoi(request.URL.Query().Get("p"))
	if err != nil {
		page = 0
	}

	if !models.ExistsUser(userId) {
		log.Debug(err)

		showError(document, request, "ユーザが存在しません。")
		return
	}

	var programs []models.Program
	i, err := models.GetProgramListByUser(sortKey, &programs, userId, true, page*10, 10)

	if err != nil {
		log.Fatal(err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	maxPage := i / 10
	if i%10 == 0 {
		maxPage--
	}

	userName, err := models.GetUserName(userId)

	if err != nil {
		log.Fatal(err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	err = tmpl.Render(document, userProgramsMember{
		DefaultMember: &templates.DefaultMember{
			Title:  userName + " - " + config.SiteTitle,
			UserID: getSessionUser(request),
		},
		Programs:     programs,
		ProgramCount: i,
		CurPage:      page,
		MaxPage:      maxPage,
		Sort:         sort,
		UserName:     userName,
		UserId:       userId,
	})
}

type userSettingsMember struct {
	*templates.DefaultMember
	UserInfo models.User
}

func userSettingsHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "userSettings.tmpl"

	userId := getSessionUser(request)

	if userId == 0 {
		log.DebugStr("匿名の管理画面へのアクセス")

		showError(document, request, "ログインが必要です。")
		// TODO: ログインさせる。
		return
	}

	var user models.User
	err := user.Load(userId)

	if err != nil {
		log.Fatal(err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	err = tmpl.Render(document, userSettingsMember{
		DefaultMember: &templates.DefaultMember{
			Title:  "管理画面 - " + config.SiteTitle,
			UserID: userId,
		},
		UserInfo: user,
	})
}
