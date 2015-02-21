package controllers

import (
	"net/http"

	"github.com/gorilla/sessions"

	"hsproom/config"
	"hsproom/templates"
)

var Router Routes
var sessionStore = sessions.NewCookieStore([]byte(config.SessionKey))

func init() {

	apiInit()

	Router.Register("/", indexHandler)
	Router.Register("/error/", flashHandler)
	Router.Register("/success/", flashHandler)
	Router.Register("/program/", programHandler)
	Router.Register("/program/list/", programListHandler)
	Router.Register("/program/view/", programViewHandler)
	Router.Register("/program/edit/", programEditHandler)
	Router.Register("/program/create/", programCreateHandler)
	Router.Register("/program/search/", programSearchHandler)
	Router.Register("/program/ranking/daily/", programRankingDailyHandler)
	Router.Register("/program/ranking/weekly/", programRankingWeeklyHandler)
	Router.Register("/program/ranking/monthly/", programRankingMonthlyHandler)
	Router.Register("/program/ranking/alltime/", programRankingAllTimeHandler)
	Router.Register("/user/logout/", userLogoutHandler)
	Router.Register("/user/view/", userViewHandler)
	Router.Register("/user/edit/", userEditHandler)
	Router.Register("/user/list/", userListHandler)
	Router.Register("/api/", apiHandler)
	Router.Register("/api/markdown/", apiMarkdownHandler)
	Router.Register("/api/twitter/search/", apiTwitterSearchHandler)
	Router.Register("/api/program/good/", apiProgramGoodHandler)
	Router.Register("/api/program/update/", apiProgramUpdateHandler)
	Router.Register("/api/program/create/", apiProgramCreateHandler)
	Router.Register("/api/program/data/", apiProgramDataHandler)
	Router.Register("/api/program/data_list/", apiProgramDataListHandler)
	Router.Register("/api/program/thumbnail/", apiProgramThumbnailHandler)
	Router.Register("/api/twitter/request_token/", apiTwitterRequestTokenHandler)
	Router.Register("/api/twitter/access_token/", apiTwitterAccessTokenHandler)
	Router.Register("/api/user/info/", apiUserInfoHandler)

}
func Del() {
	apiDel()
}

func getSessionUser(request *http.Request) int {
	session, _ := sessionStore.Get(request, "go-wiki")
	if session.Values["User"] == nil {
		return 0
	}

	return session.Values["User"].(int)
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
	session, _ := sessionStore.Get(request, "user")

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

	session, _ := sessionStore.Get(request, "user")

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
			Title: message,
			User:  getSessionUser(request),
		},
		Messages: flashes,
		Referer:  request.Referer(),
	})
}
