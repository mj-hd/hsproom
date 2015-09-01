package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"

	"../config"
	"../templates"
	"../utils/log"
)

var Router Routes
var sessionStore = sessions.NewCookieStore([]byte(config.SessionKey))

func Init() {

	apiInit()

	Router.RegisterPage("/", indexHandler)
	Router.Register("/error/", flashHandler)
	Router.Register("/success/", flashHandler)
	Router.RegisterPage("/help/", helpHandler)
	Router.RegisterPage("/about/", aboutHandler)
	Router.RegisterPage("/program/", programHandler)
	Router.RegisterPage("/program/list/", programListHandler)
	Router.RegisterPage("/program/view/", programViewHandler)
	Router.RegisterPage("/program/edit/", programEditHandler)
	Router.RegisterPage("/program/create/", programCreateHandler)
	Router.RegisterPage("/program/search/", programSearchHandler)
	Router.Register("/program/remote_view/", programRemoteViewHandler)
	Router.RegisterPage("/program/ranking/daily/", programRankingDailyHandler)
	Router.RegisterPage("/program/ranking/weekly/", programRankingWeeklyHandler)
	Router.RegisterPage("/program/ranking/monthly/", programRankingMonthlyHandler)
	Router.RegisterPage("/program/ranking/alltime/", programRankingAllTimeHandler)
	Router.RegisterPage("/source/create/", sourceCreateHandler)
	Router.RegisterPage("/source/edit/", sourceEditHandler)
	Router.RegisterPage("/user/logout/", userLogoutHandler)
	Router.RegisterPage("/user/login/", userLoginHandler)
	Router.RegisterPage("/user/view/", userViewHandler)
	Router.RegisterPage("/user/edit/", userEditHandler)
	Router.RegisterPage("/user/list/", userListHandler)
	Router.RegisterPage("/user/programs/", userProgramsHandler)
	Router.RegisterPage("/user/settings/", userSettingsHandler)
	Router.RegisterApi("/api/", apiHandler)
	Router.RegisterGetApi("/api/markdown/", apiMarkdownHandler)
	Router.RegisterGetApi("/api/twitter/search/", apiTwitterSearchHandler)
	Router.RegisterPostApi("/api/program/good/", apiProgramGoodHandler)
	Router.RegisterGetApi("/api/program/good/count/", apiProgramGoodCountHandler)
	Router.RegisterPostApi("/api/program/update/", apiProgramUpdateHandler)
	Router.RegisterPostApi("/api/program/create/", apiProgramCreateHandler)
	Router.RegisterPostApi("/api/program/remove/", apiProgramRemoveHandler)
	Router.Register("/api/program/data/", apiProgramDataHandler)
	Router.RegisterGetApi("/api/program/data_list/", apiProgramDataListHandler)
	Router.Register("/api/program/thumbnail/", apiProgramThumbnailHandler)
	Router.RegisterGetApi("/api/twitter/request_token/", apiTwitterRequestTokenHandler)
	Router.Register("/api/twitter/access_token/", apiTwitterAccessTokenHandler)
	Router.RegisterGetApi("/api/google/request_token/", apiGoogleRequestTokenHandler)
	Router.Register("/api/google/access_token/", apiGoogleAccessTokenHandler)
	Router.RegisterGetApi("/api/user/info/", apiUserInfoHandler)
	Router.RegisterGetApi("/api/user/programs/", apiUserProgramsHandler)
	Router.RegisterGetApi("/api/user/goods/", apiUserGoodsHandler)
	Router.RegisterPostApi("/api/good/remove/", apiGoodRemoveHandler)
	Router.RegisterGetApi("/api/comment/list/", apiCommentListHandler)
	Router.RegisterPostApi("/api/comment/post/", apiCommentPostHandler)
	Router.RegisterPostApi("/api/comment/delete/", apiCommentDeleteHandler)

}
func Del() {
	apiDel()
}

func getSession(request *http.Request) (*sessions.Session, error) {
	session, err := sessionStore.Get(request, config.SessionName)

	if err == nil {
		// 一週間
		session.Options.MaxAge = 86400 * 7
	}

	return session, err
}

func getSessionUser(request *http.Request) int {
	session, _ := getSession(request)
	if session.Values["User"] == nil {
		return 0
	}

	return session.Values["User"].(int)
}

func removeSession(document http.ResponseWriter, request *http.Request) {
	session, err := getSession(request)
	if err != nil {
		return
	}

	session.Options = &sessions.Options{MaxAge: -1, Path: "/"}
	session.Save(request, document)
}

func writeStruct(document http.ResponseWriter, s interface{}, httpStatus int) {

	var err error

	document.Header().Set("Content-Type", "application/json")
	jso, err := json.Marshal(s)

	if err != nil {

		log.Fatal(err)

		document.WriteHeader(500)
		document.Write([]byte("{ \"Status\" : \"error\", \"Message\" : \"不明のエラーです。\" }"))

		return
	}

	document.WriteHeader(httpStatus)
	document.Write(jso)
}

type apiMember struct {
	Status  string
	Message string
}

type Routes struct {
	keys   []string
	values []func(http.ResponseWriter, *http.Request)
}
type Route struct {
	Path     string
	Function func(http.ResponseWriter, *http.Request)
}

func (this *Routes) Iterator() <-chan Route {
	ret := make(chan Route)

	go func() {
		for i, k := range this.keys {
			var route Route
			route.Path = k
			route.Function = this.values[i]

			ret <- route
		}
		close(ret)
	}()

	return ret
}

func (this *Routes) RegisterPage(path string, fn func(http.ResponseWriter, *http.Request) error) {
	this.keys = append(this.keys, path)
	this.values = append(this.values, func(document http.ResponseWriter, request *http.Request) {
		err := fn(document, request)
		if err != nil {
			log.FatalStr("ページの表示に失敗:")
			log.Fatal(err)

			showError(document, request, "ページの表示中にエラーが発生しました。管理人へ報告してください。")
		}
	})
}

func (this *Routes) RegisterApi(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	this.keys = append(this.keys, path)
	this.values = append(this.values, func(document http.ResponseWriter, request *http.Request) {
		status, err := fn(document, request)
		if err != nil {
			log.FatalStr("APIの実行に失敗:")
			log.Fatal(err)

			writeStruct(document, apiMember{
				Status:  "error",
				Message: err.Error(),
			}, status)

			return
		}
	})
}

func (this *Routes) RegisterPostApi(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	wrapper := func(document http.ResponseWriter, request *http.Request) (int, error) {
		if request.Method != "POST" {
			return http.StatusBadRequest, errors.New("POST以外のメソッド")
		}

		return fn(document, request)
	}

	this.RegisterApi(path, wrapper)
}

func (this *Routes) RegisterGetApi(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	wrapper := func(document http.ResponseWriter, request *http.Request) (int, error) {
		if request.Method != "GET" {
			return http.StatusBadRequest, errors.New("GET以外のメソッド")
		}

		return fn(document, request)
	}

	this.RegisterApi(path, wrapper)
}

func (this *Routes) Register(path string, fn func(http.ResponseWriter, *http.Request)) {
	this.keys = append(this.keys, path)
	this.values = append(this.values, fn)
}

func (this *Routes) Value(path string) func(http.ResponseWriter, *http.Request) {
	for i, key := range this.keys {
		if path == key {
			return this.values[i]
		}
	}
	return nil
}

func (this *Routes) Key(fn *func(http.ResponseWriter, *http.Request)) string {
	for i, val := range this.values {
		if fn == &val {
			return this.keys[i]
		}
	}
	return ""
}

func showError(document http.ResponseWriter, request *http.Request, message string) {
	session, _ := getSession(request)

	session.AddFlash(message)
	session.Save(request, document)
	http.Redirect(document, request, "/error/", http.StatusSeeOther)
}

type flashMember struct {
	*templates.DefaultMember
	Messages []interface{}
	Referer  string
}

func flashHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "flash.tmpl"

	session, _ := getSession(request)

	var message string
	if request.URL.Path == "/error/" {
		message = "エラー"
	} else {
		message = "成功"
	}

	flashes := session.Flashes()
	session.Save(request, document)

	tmpl.Render(document, flashMember{
		DefaultMember: &templates.DefaultMember{
			Title:  message,
			UserID: getSessionUser(request),
		},
		Messages: flashes,
		Referer:  request.Referer(),
	})
}
