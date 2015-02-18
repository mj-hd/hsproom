package controllers

import (
	"hsproom/config"
	"hsproom/models"
	"hsproom/templates"
	"hsproom/utils"
	"net/http"
	"os"
	"strconv"
)

type userViewMember struct {
	*templates.DefaultMember
	UserInfo     *models.User
	UserPrograms *[]models.ProgramInfo
}

func userViewHandler(document http.ResponseWriter, request *http.Request) {

	var err error

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "userView.tmpl"

	rawUid := request.URL.Query().Get("u")
	uid, err := strconv.Atoi(rawUid)

	if err != nil {
		utils.PromulgateDebug(os.Stdout, err)

		showError(document, request, "ユーザが見つかりませんでした。")
		return
	}

	var user models.User
	err = user.Load(uid)

	if err != nil {
		utils.PromulgateDebug(os.Stdout, err)

		showError(document, request, "ユーザが見つかりませんでした。")
		return
	}

	var programs []models.ProgramInfo

	_, err = models.GetProgramListByUser(&programs, user.Name, 4, 0)

	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	err = tmpl.Render(document, userViewMember{
		DefaultMember: &templates.DefaultMember{
			Title: user.Name + " のプロフィール - " + config.SiteTitle,
			User:  getSessionUser(request),
		},
		UserInfo:     &user,
		UserPrograms: &programs,
	})
	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)

		showError(document, request, "ページの表示に失敗しました。管理人へ問い合わせてください。")
	}
}

func userListHandler(document http.ResponseWriter, request *http.Request) {

}

func userLogoutHandler(document http.ResponseWriter, request *http.Request) {

	session, err := sessionStore.Get(request, "go-wiki")

	if err != nil {
		utils.PromulgateDebug(os.Stdout, err)

		showError(document, request, "ログアウトに失敗しました。")

		return
	}

	session.Values["User"] = nil
	session.Save(request, document)

	http.Redirect(document, request, "http://localhost:8080/", 301)
}

func userEditHandler(document http.ResponseWriter, request *http.Request) {

}