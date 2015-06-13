package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"../config"
	"../models"
	"../bot"
	"../utils/log"
	"../utils/twitter"
	"../utils/google"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

var twitterClient *twitter.Client
var oauthClient *twitter.OAuthClient
var oauth2Client *google.OAuth2Client

func apiInit() {
	var err error

	// 2-legged
	twitterClient, err = twitter.NewClient(config.TwitterKey, config.TwitterSecret)
	if err != nil {
		panic(err)
	}

	// 3-legged
	oauthClient, err = twitter.NewOAuthClient(config.TwitterKey, config.TwitterSecret)
	if err != nil {
		panic(err)
	}

	oauth2Client, err = google.NewOAuth2Client(config.GoogleKey, config.GoogleSecret)
	if err != nil {
		panic(err)
	}
}
func apiDel() {
	twitterClient.Close()
}

func writeStruct(document http.ResponseWriter, s interface{}, httpStatus int) {

	var err error

	document.Header().Set("Content-Type", "application/json")
	jso, err := json.Marshal(s)

	if err != nil {

		log.Fatal(os.Stdout, err)

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

func apiHandler(document http.ResponseWriter, request *http.Request) {

	writeStruct(document, apiMember{
		Status:  "success",
		Message: "エラーなし",
	}, 200)

}

type apiProgramGoodMember struct {
	*apiMember
}

func apiProgramGoodHandler(document http.ResponseWriter, request *http.Request) {

	var err error

	if request.Method != "POST" {

		log.DebugStr(os.Stdout, "POST以外のGoodリクエスト")

		writeStruct(document, apiProgramGoodMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "POST以外のメソッドです。",
			},
		}, 400)

		return
	}

	// プログラムIDの取得
	rawProgramId := request.FormValue("p")
	programId, err := strconv.Atoi(rawProgramId)

	if err != nil {

		log.DebugStr(os.Stdout, "不正なプログラムID "+string(programId))
		log.Debug(os.Stdout, err)

		writeStruct(document, apiProgramGoodMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "プログラムIDが不正です。",
			},
		}, 400)

		return
	}

	if !models.ExistsProgram(programId) {

		log.FatalStr(os.Stdout, "プログラム["+string(programId)+"]の読み込みに失敗")
		log.Fatal(os.Stdout, err)

		writeStruct(document, apiProgramGoodMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "内部エラーが発生しました。",
			},
		}, 500)

		return
	}

	user := getSessionUser(request)
	if user == 0 {

		log.DebugStr(os.Stdout, "匿名のGoodリクエスト")

		writeStruct(document, apiProgramGoodMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "ログインが必要です。",
			},
		}, 400)

		return
	}

	if !models.CanGoodProgram(user, programId) {
		log.DebugStr(os.Stdout, "無効なGoodリクエスト")

		writeStruct(document, apiProgramGoodMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "無効なGoodリクエストです。",
			},
		}, 400)

		return
	}

	var good models.Good
	good.User = user
	good.Program = programId

	_, err = good.Create()
	if err != nil {

		log.Fatal(os.Stdout, err)
		writeStruct(document, apiProgramGoodMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "内部エラーが発生しました。",
			},
		}, 500)

		return
	}

	writeStruct(document, apiProgramGoodMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "保存に成功しました。",
		},
	}, 200)
}

type apiProgramUpdateMember struct {
	*apiMember
}

func apiProgramUpdateHandler(document http.ResponseWriter, request *http.Request) {

	// メソッドの確認
	if request.Method != "POST" {

		log.DebugStr(os.Stdout, "POST以外のUpdateリクエスト")

		writeStruct(document, apiProgramUpdateMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "POST以外のメソッドです。",
			},
		}, 400)

		return
	}

	// 入力値のバリデート
	var rawProgram models.RawProgram
	targetFlags := models.ProgramId | models.ProgramTitle | models.ProgramThumbnail | models.ProgramDescription | models.ProgramStartax | models.ProgramAttachments | models.ProgramSteps | models.ProgramSourcecode | models.ProgramRuntime

	rawProgram.Id = request.FormValue("id")
	rawProgram.Title = bluemonday.UGCPolicy().Sanitize(request.FormValue("title"))
	rawProgram.Thumbnail = request.FormValue("thumbnail")
	rawProgram.Description = request.FormValue("description")
	rawProgram.Startax = request.FormValue("startax")
	rawProgram.Attachments = request.FormValue("attachments")
	rawProgram.Steps = request.FormValue("steps")
	rawProgram.Sourcecode = request.FormValue("sourcecode")
	rawProgram.Runtime = request.FormValue("runtime")

	if rawProgram.Steps == "" {
		targetFlags -= models.ProgramSteps
	}
	if rawProgram.Sourcecode == "" {
		targetFlags -= models.ProgramSourcecode
	}

	err := rawProgram.Validate(targetFlags)
	if err != nil {
		log.Debug(os.Stdout, err)

		writeStruct(document, apiProgramUpdateMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: err.Error(),
			},
		}, 400)

		return
	}

	// プログラムへ変換
	program, err := rawProgram.ToProgram(targetFlags)
	if err != nil {
		log.Debug(os.Stdout, err)

		writeStruct(document, apiProgramUpdateMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "内部エラーが発生しました。",
			},
		}, 500)

		return
	}

	// プログラムの確認
	var prevProgInfo models.ProgramInfo

	err = prevProgInfo.Load(program.Id)
	if err != nil {
		log.Debug(os.Stdout, err)

		writeStruct(document, apiProgramUpdateMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "プログラムIDが不正です。",
			},
		}, 400)

		return
	}

	// ユーザのチェック
	if getSessionUser(request) != prevProgInfo.User {
		log.DebugStr(os.Stdout, "プログラムの権限のない変更")

		writeStruct(document, apiProgramUpdateMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "プログラムを編集する権限がありません。",
			},
		}, 400)

		return
	}

	// 適用
	prevProgInfo.Title = program.Title
	prevProgInfo.Description = program.Description
	prevProgInfo.Steps = program.Steps
	prevProgInfo.Runtime = program.Runtime

	// 以前のプログラムと合成する
	program.ProgramInfo = &prevProgInfo

	err = program.Update()
	if err != nil {
		log.Fatal(os.Stdout, err)

		writeStruct(document, apiProgramUpdateMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "保存に失敗しました。",
			},
		}, 500)

		return
	}

	writeStruct(document, apiProgramUpdateMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "保存に成功しました。",
		},
	}, 200)
}

type apiProgramCreateMember struct {
	*apiMember
	Id int
}

type apiNameValuePair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func apiProgramCreateHandler(document http.ResponseWriter, request *http.Request) {

	// メソッドの確認
	if request.Method != "POST" {

		log.DebugStr(os.Stdout, "POST以外のCreateリクエスト")

		writeStruct(document, apiProgramUpdateMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "POST以外のメソッドです。",
			},
		}, 400)

		return
	}

	// 入力値のバリデート
	var rawProgram models.RawProgram
	targetFlags := models.ProgramTitle | models.ProgramThumbnail | models.ProgramDescription | models.ProgramStartax | models.ProgramAttachments | models.ProgramSteps | models.ProgramSourcecode | models.ProgramRuntime

	rawProgram.Title = bluemonday.UGCPolicy().Sanitize(request.FormValue("title"))
	rawProgram.Thumbnail = request.FormValue("thumbnail")
	rawProgram.Description = request.FormValue("description")
	rawProgram.Startax = request.FormValue("startax")
	rawProgram.Attachments = request.FormValue("attachments")
	rawProgram.Steps = request.FormValue("steps")
	rawProgram.Sourcecode = request.FormValue("sourcecode")
	rawProgram.Runtime = request.FormValue("runtime")

	if rawProgram.Steps == "" {
		targetFlags -= models.ProgramSteps
	}
	if rawProgram.Sourcecode == "" {
		targetFlags -= models.ProgramSourcecode
	}

	userId := getSessionUser(request)

	// ログインしていない
	if userId == 0 {
		log.DebugStr(os.Stdout, "匿名のCreateリクエスト")

		writeStruct(document, apiProgramCreateMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "ログインする必要があります。",
			},
		}, 400)

		return
	}

	err := rawProgram.Validate(targetFlags)
	if err != nil {
		log.Debug(os.Stdout, err)

		writeStruct(document, apiProgramCreateMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: err.Error(),
			},
		}, 400)

		return
	}

	// プログラムへ変換
	program, err := rawProgram.ToProgram(targetFlags)
	if err != nil {
		log.Debug(os.Stdout, err)

		writeStruct(document, apiProgramCreateMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "内部エラーが発生しました。",
			},
		}, 500)

		return
	}

	program.User = userId

	id, err := program.Create()
	if err != nil {
		log.Fatal(os.Stdout, err)

		writeStruct(document, apiProgramCreateMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "保存に失敗しました。もう一度お試しください。",
			},
		}, 500)

		return
	}

	writeStruct(document, apiProgramCreateMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "保存に成功しました。",
		},
		Id: id,
	}, 200)

	bot.UpdateTweet("新しいプログラムが投稿されました! #hsproom\n\n "+ program.Title + " by "+program.UserName + " " +config.SiteURL+"/program/view/?p="+strconv.Itoa(id))
}

type apiProgramDataListMember struct {
	*apiMember
	Names []string
}

func apiProgramDataListHandler(document http.ResponseWriter, request *http.Request) {

	if request.Method != "GET" {

		log.DebugStr(os.Stdout, "GET以外のDataListリクエスト")

		writeStruct(document, apiProgramDataListMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "GETを使用してください。",
			},
		}, 400)

		return
	}

	programId, err := strconv.Atoi(request.URL.Query().Get("p"))

	if err != nil {

		log.DebugStr(os.Stdout, "プログラムIDが不正")

		writeStruct(document, apiProgramDataListMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "プログラムIDが不正です。",
			},
		}, 400)

		return
	}

	program := models.NewProgram()
	err = program.Load(programId)

	if err != nil {

		log.Debug(os.Stdout, err)

		writeStruct(document, apiProgramDataListMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "プログラムが存在しません。",
			},
		}, 400)

		return
	}

	var names []string

	for _, file := range program.Attachments.Files {
		names = append(names, file.Name)
	}

	writeStruct(document, apiProgramDataListMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "添付ファイル一覧の取得に成功しました。",
		},
		Names: names,
	}, 200)
}

// jsonじゃないよ
func apiProgramDataHandler(document http.ResponseWriter, request *http.Request) {

	if request.Method != "GET" {

		log.DebugStr(os.Stdout, "GET以外のDataリクエスト")

		document.WriteHeader(400)

		return
	}

	document.Header().Set("Content-Type", "application/octet-stream")

	rawProgramId := request.URL.Query().Get("p")
	programId, err := strconv.Atoi(rawProgramId)

	if err != nil {

		log.DebugStr(os.Stdout, "不正なプログラムID")

		document.WriteHeader(400)

		return
	}

	program := models.NewProgram()
	err = program.Load(programId)

	if err != nil {

		log.Debug(os.Stdout, err)
		log.DebugStr(os.Stdout, "プログラムが見つからない")

		document.WriteHeader(404)

		return
	}

	fileName := request.URL.Query().Get("f")
	if fileName == "" {

		log.DebugStr(os.Stdout, "空のDataリクエスト")

		document.WriteHeader(404)

		return
	}

	if fileName == "start.ax" {

		document.WriteHeader(200)
		document.Write(program.Startax)

		return
	}

	// ファイルを検索する

	for _, file := range program.Attachments.Files {
		if file.Name == fileName {

			document.WriteHeader(200)
			document.Write(file.Data)

			return
		}
	}

	// ファイルが見つからなかった
	log.DebugStr(os.Stdout, "指定されたファイルが見つからなかった")
	document.WriteHeader(404)
}

// jsonじゃない
func apiProgramThumbnailHandler(document http.ResponseWriter, request *http.Request) {

	if request.Method != "GET" {

		log.DebugStr(os.Stdout, "GET以外のThumbnailリクエスト")

		document.WriteHeader(400)

		return
	}

	document.Header().Set("Content-Type", "image/png")

	rawProgramId := request.URL.Query().Get("p")
	programId, err := strconv.Atoi(rawProgramId)

	if err != nil {

		log.Debug(os.Stdout, err)

		document.WriteHeader(400)

		return
	}

	program := models.NewProgram()

	err = program.Load(programId)

	if err != nil {

		log.Debug(os.Stdout, err)

		document.WriteHeader(404)

		return
	}

	document.Write(program.Thumbnail)
}

type apiMarkdownMember struct {
	*apiMember
	Markdown string
}

func apiMarkdownHandler(document http.ResponseWriter, request *http.Request) {

	var text string
	var texts = request.URL.Query()["text"]

	if len(texts) == 0 {
		text = ""
	} else {
		text = texts[0]
	}

	writeStruct(document, apiMarkdownMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "変換に成功しました。",
		},
		Markdown: string(
			bluemonday.UGCPolicy().SanitizeBytes(
				blackfriday.MarkdownCommon(
					[]byte(text))))}, 200)

}

type apiTwitterSearchMember struct {
	*apiMember
	Tweets twitter.SearchResponse
}

func apiTwitterSearchHandler(document http.ResponseWriter, request *http.Request) {

	programName := bluemonday.UGCPolicy().Sanitize(request.URL.Query().Get("pn"))
	rawNumber := request.URL.Query().Get("n")
	rawOffset := request.URL.Query().Get("o")

	number, err := strconv.Atoi(rawNumber)
	if err != nil {
		log.Debug(os.Stdout, err)

		writeStruct(document, apiTwitterSearchMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "不正な数値です。",
			},
		}, 403)
		return
	}

	offset, err := strconv.ParseInt(rawOffset, 10, 64)
	if err != nil {
		log.Debug(os.Stdout, err)

		writeStruct(document, apiTwitterSearchMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "不正な数値です。",
			},
		}, 403)
		return
	}

	query := "#hsproom"

	symbols, err := regexp.Compile("[\\s\x21-\x2f\x3a-\x40\x5b-\x60\x7b-\x7e]")

	if err != nil {
		log.Fatal(os.Stdout, err)

		writeStruct(document, apiTwitterSearchMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "検索に失敗しました。",
			},
		}, 500)
		return
	}

	programName = symbols.ReplaceAllString(programName, "")

	if programName != "" {
		query += " #" + programName
	}

	tweets, err := twitterClient.SearchTweets(query, number, offset)
	if err != nil {
		log.Fatal(os.Stdout, err)

		writeStruct(document, apiTwitterSearchMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "検索に失敗しました。",
			},
		}, 500)
		return
	}

	if err != nil {
		log.Fatal(os.Stdout, err)

		writeStruct(document, apiTwitterSearchMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "内部エラーが発生しました。",
			},
		}, 500)
		return
	}

	writeStruct(document, apiTwitterSearchMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "検索に成功しました。",
		},
		Tweets: tweets,
	}, 200)

}

type apiOAuthRequestTokenMember struct {
	*apiMember
	AuthURL string
}

func apiTwitterRequestTokenHandler(document http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		log.DebugStr(os.Stdout, "GET以外のRequestTokenリクエスト")

		writeStruct(document, apiOAuthRequestTokenMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "GETを使用してください。",
			},
		}, 403)
		return
	}

	callbackUrl := request.URL.Query().Get("c")

	if callbackUrl == "" {
		log.DebugStr(os.Stdout, "callback指定のないRequestTokenリクエスト")

		writeStruct(document, apiOAuthRequestTokenMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "コールバック指定が必要です。",
			},
		}, 400)
		return
	}

	url, err := oauthClient.GetAuthURL(config.SiteURL + "/api/twitter/access_token/?c=" + url.QueryEscape(callbackUrl))
	if err != nil {
		log.Fatal(os.Stdout, err)

		writeStruct(document, apiOAuthRequestTokenMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "内部エラーが発生しました。",
			},
		}, 500)
		return
	}

	writeStruct(document, apiOAuthRequestTokenMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "URLの取得に成功しました。",
		},
		AuthURL: url,
	}, 200)
}

// jsonじゃない
func apiTwitterAccessTokenHandler(document http.ResponseWriter, request *http.Request) {

	if request.Method != "GET" {
		log.DebugStr(os.Stdout, "GET以外のAccessTokenリクエスト")

		document.WriteHeader(403)

		return
	}

	verifier := request.URL.Query().Get("oauth_verifier")
	token := request.URL.Query().Get("oauth_token")

	if verifier == "" || token == "" {
		log.DebugStr(os.Stdout, "クエリが空")

		document.WriteHeader(403)

		return
	}

	accessToken, err := oauthClient.GetAccessToken(verifier, token)
	if err != nil {
		log.Fatal(os.Stdout, err)

		document.WriteHeader(500)

		return
	}

	user, err := oauthClient.CheckUserCredentialsAndGetUser(accessToken)

	if err != nil {
		log.Fatal(os.Stdout, err)

		document.WriteHeader(500)

		return
	}

	var dbUser models.User
	dbUser.ScreenName = user.ScreenName
	dbUser.Name = user.Name
	dbUser.Profile = user.Description
	dbUser.IconURL = user.ProfileImageURL
	dbUser.Website = "https://twitter.com/" + user.ScreenName
	dbUser.Location = user.Location

	var id int

	var oldUser models.User
	err = oldUser.LoadFromScreenName(user.ScreenName)

	if err != nil {
		id, err = dbUser.Create()
		if err != nil {
			log.Fatal(os.Stdout, err)

			document.WriteHeader(500)

			return
		}
	} else {

		id = oldUser.Id
		dbUser.Id = oldUser.Id

		err = dbUser.Update()

		if err != nil {
			log.Fatal(os.Stdout, err)

			document.WriteHeader(500)

			return
		}

	}

	session, err := getSession(request)
	if err != nil {
		log.Fatal(os.Stdout, err)

		document.WriteHeader(403)

		return
	}

	session.Values["User"] = id
	session.Save(request, document)

	callbackUrl := request.URL.Query().Get("c")
	if callbackUrl == "" {
		callbackUrl = config.SiteURL + "/"
	}

	http.Redirect(document, request, callbackUrl, 301)
	return
}

func apiGoogleRequestTokenHandler(document http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		log.DebugStr(os.Stdout, "GET以外のGoogleRequestTokenリクエスト")

		writeStruct(document, apiOAuthRequestTokenMember{
			apiMember: &apiMember{
				Status: "error",
				Message: "GETを使用してください。",
			},
		}, 403)
		return
	}

	url, err := oauth2Client.GetAuthURL(config.SiteURL + "/api/google/access_token/")
	if err != nil {
		log.Fatal(os.Stdout, err)

		writeStruct(document, apiOAuthRequestTokenMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "内部エラーが発生しました。",
			},
		}, 500)
		return
	}

	session, err := getSession(request)
	if err != nil {
		log.Debug(os.Stdout, err)

		writeStruct(document, apiOAuthRequestTokenMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "クッキーが有効ではありません。",
			},
		}, 400)
		return
	}

	callbackUrl := request.URL.Query().Get("c")
	if callbackUrl == "" {
		callbackUrl = config.SiteURL
	}

	session.Values["Callback"] = callbackUrl
	session.Save(request, document)

	writeStruct(document, apiOAuthRequestTokenMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "URLの取得に成功しました。",
		},
		AuthURL: url,
	}, 200)
}

func apiGoogleAccessTokenHandler(document http.ResponseWriter, request *http.Request) {

	if request.Method != "GET" {
		log.DebugStr(os.Stdout, "GET以外のGoogleAccessTokenリクエスト")

		document.WriteHeader(403)

		return
	}

	verifier := request.URL.Query().Get("state")
	token := request.URL.Query().Get("code")

	if verifier == "" || token == "" {
		log.DebugStr(os.Stdout, "クエリが空")

		document.WriteHeader(403)

		return
	}

	accessToken, err := oauth2Client.GetToken(verifier, token)
	if err != nil {
		log.Fatal(os.Stdout, err)

		document.WriteHeader(500)

		return
	}

	userinfo, err := oauth2Client.GetUser(accessToken)

	if err != nil {
		log.Fatal(os.Stdout, err)

		document.WriteHeader(500)

		return
	}

	var dbUser models.User
	dbUser.ScreenName = userinfo.IdString
	dbUser.Name       = userinfo.Name
	dbUser.Profile    = ""
	dbUser.IconURL    = userinfo.Picture
	dbUser.Location   = userinfo.Locale
	dbUser.Website    = userinfo.Link

	var id int

	var oldUser models.User
	err = oldUser.LoadFromScreenName(userinfo.IdString)

	if err != nil {
		id, err = dbUser.Create()

		if err != nil {
			log.Fatal(os.Stdout, err)

			document.WriteHeader(500)

			return
		}
	} else {
	
		id = oldUser.Id
		dbUser.Id = oldUser.Id

		err = dbUser.Update()

		if err != nil {
			log.Fatal(os.Stdout, err)

			document.WriteHeader(500)

			return
		}
	}

	session, err := getSession(request)
	if err != nil {
		log.Fatal(os.Stdout, err)

		document.WriteHeader(403)

		return
	}

	session.Values["User"] = id

	callbackUrl := session.Values["Callback"].(string)

	session.Values["Callback"] = ""
	session.Save(request, document)

	if callbackUrl == "" {
		callbackUrl = config.SiteURL + "/"
	}

	http.Redirect(document, request, callbackUrl, 301)
	return
}

type apiUserInfoMember struct {
	*apiMember
	*models.User
}

func apiUserInfoHandler(document http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		log.DebugStr(os.Stdout, "GET以外のUserInfoリクエスト")

		writeStruct(document, apiUserInfoMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "GETを使用してください。",
			},
		}, 400)
		return
	}

	rawUserId := request.URL.Query().Get("u")
	userId, err := strconv.Atoi(rawUserId)

	if err != nil {
		log.DebugStr(os.Stdout, "ユーザ指定のないリクエスト")

		writeStruct(document, apiUserInfoMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "数値が不正です。",
			},
		}, 400)

		return
	}

	var user models.User
	err = user.Load(userId)

	if err != nil {
		log.DebugStr(os.Stdout, "存在しないユーザのリクエスト")

		writeStruct(document, apiUserInfoMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "ユーザIDが不正です。",
			},
		}, 400)

		return
	}

	writeStruct(document, user, 200)

}

type apiUserProgramListMember struct {
	*apiMember
	Programs     []models.ProgramInfo
	ProgramCount int
}

func apiUserProgramsHandler(document http.ResponseWriter, request *http.Request) {

	if request.Method != "GET" {
		log.DebugStr(os.Stdout, "GET以外のUserProgramListリクエスト")

		writeStruct(document, apiUserProgramListMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "GETを使用してください。",
			},
		}, 403)
		return
	}

	userId, err := strconv.Atoi(request.URL.Query().Get("u"))
	if err != nil {
		log.Debug(os.Stdout, err)

		writeStruct(document, apiUserProgramListMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "ユーザIDが不正です。",
			},
		}, 403)
		return
	}

	if !models.ExistsUser(userId) {
		log.Fatal(os.Stdout, err)

		writeStruct(document, apiUserProgramListMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "ユーザが存在しません。",
			},
		}, 403)
		return
	}

	offset, err := strconv.Atoi(request.URL.Query().Get("o"))
	if err != nil {
		offset = 0
	}

	number, err := strconv.Atoi(request.URL.Query().Get("n"))
	if err != nil {
		log.Debug(os.Stdout, err)

		writeStruct(document, apiUserProgramListMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "取得数の指定が不正です。",
			},
		}, 403)
	}

	var programs []models.ProgramInfo
	i, err := models.GetProgramListByUser(models.ProgramColCreated, &programs, userId, true, offset, number)

	if err != nil {
		log.Fatal(os.Stdout, err)

		writeStruct(document, apiUserProgramListMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "内部エラーが発生しました。",
			},
		}, 500)
		return
	}

	writeStruct(document, apiUserProgramListMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "一覧の取得に成功しました。",
		},
		Programs:     programs,
		ProgramCount: i,
	}, 200)
}

type apiUserGoodsMember struct {
	*apiMember
	Programs     []models.ProgramInfo
	ProgramCount int
}

func apiUserGoodsHandler(document http.ResponseWriter, request *http.Request) {

	if request.Method != "GET" {
		log.DebugStr(os.Stdout, "GET以外のUserGoodリクエスト")

		writeStruct(document, apiUserGoodsMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "GETを使用してください。",
			},
		}, 403)
		return
	}

	userId, err := strconv.Atoi(request.URL.Query().Get("u"))

	if err != nil {
		log.Debug(os.Stdout, err)

		writeStruct(document, apiUserGoodsMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "不正なユーザIDです。",
			},
		}, 403)
		return
	}

	if !models.ExistsUser(userId) {
		log.Debug(os.Stdout, err)

		writeStruct(document, apiUserGoodsMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "存在しないユーザIDです。",
			},
		}, 403)
		return
	}

	offset, err := strconv.Atoi(request.URL.Query().Get("o"))

	if err != nil {
		log.DebugStr(os.Stdout, "不正なオフセット。")

		writeStruct(document, apiUserGoodsMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "オフセット値が不正です。",
			},
		}, 403)
	}

	number, err := strconv.Atoi(request.URL.Query().Get("n"))

	if err != nil {
		log.DebugStr(os.Stdout, "不正な制限数。")

		writeStruct(document, apiUserGoodsMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "取得数が不正です。",
			},
		}, 403)
		return
	}

	var goods []models.Good

	programCount, err := models.GetGoodListByUser(&goods, userId, offset, number)

	if err != nil {

		log.Fatal(os.Stdout, err)

		writeStruct(document, apiUserGoodsMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "内部エラーが発生しました。",
			},
		}, 500)
		return
	}

	programs := make([]models.ProgramInfo, number)

	for i, good := range goods {

		if good.Program > 0 {
			err = programs[i].Load(good.Program)

			if err != nil {
			}
		}
	}

	writeStruct(document, apiUserGoodsMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "取得に成功しました。",
		},
		Programs:     programs,
		ProgramCount: programCount,
	}, 200)
}

type apiProgramGoodCountMember struct {
	*apiMember
	GoodCount int
}

func apiProgramGoodCountHandler(document http.ResponseWriter, request *http.Request) {

	if request.Method != "GET" {
		log.DebugStr(os.Stdout, "GET以外のProgramGoodCountリクエスト")

		writeStruct(document, apiProgramGoodCountMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "GETを使用してください。",
			},
		}, 403)
		return
	}

	programId, err := strconv.Atoi(request.URL.Query().Get("p"))

	if err != nil {
		log.Debug(os.Stdout, err)

		writeStruct(document, apiProgramGoodCountMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "不正なプログラムIDです。",
			},
		}, 403)
		return
	}

	if !models.ExistsProgram(programId) {
		log.DebugStr(os.Stdout, "存在しないプログラムID。")

		writeStruct(document, apiProgramGoodCountMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "存在しないプログラムIDです。",
			},
		}, 403)
		return
	}

	goodCount := models.GetGoodCountByProgram(programId)

	writeStruct(document, apiProgramGoodCountMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "いいね数の取得に成功しました。",
		},
		GoodCount: goodCount,
	}, 200)

}

type apiProgramRemoveMember struct {
	*apiMember
}

func apiProgramRemoveHandler(document http.ResponseWriter, request *http.Request) {

	if request.Method != "POST" {
		log.DebugStr(os.Stdout, "POST以外のProgramRemoveリクエスト")

		writeStruct(document, apiProgramRemoveMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "POSTを使用してください。",
			},
		}, 403)
		return
	}

	programId, err := strconv.Atoi(request.FormValue("p"))

	if err != nil {
		log.Debug(os.Stdout, err)

		writeStruct(document, apiProgramRemoveMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "不正なプログラムIDです。",
			},
		}, 403)
		return
	}

	userId := getSessionUser(request)
	if userId == 0 {
		log.DebugStr(os.Stdout, "匿名のProgramRemoveリクエスト")

		writeStruct(document, apiProgramRemoveMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "削除する権限がありません。",
			},
		}, 403)
		return
	}

	var program models.ProgramInfo
	err = program.Load(programId)

	if err != nil {
		log.Fatal(os.Stdout, err)

		writeStruct(document, apiProgramRemoveMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "内部エラーが発生しました。",
			},
		}, 500)
		return
	}

	if program.User != userId {
		log.DebugStr(os.Stdout, "権限のないProgramRemoveリクエスト")

		writeStruct(document, apiProgramRemoveMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "削除する権限がありません。",
			},
		}, 403)
		return
	}

	err = program.Remove()

	if err != nil {
		log.Fatal(os.Stdout, err)

		writeStruct(document, apiProgramRemoveMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "内部エラーが発生しました。",
			},
		}, 500)
		return
	}

	writeStruct(document, apiProgramRemoveMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "削除に成功しました。",
		},
	}, 200)
}

type apiGoodRemoveMember struct {
	*apiMember
}

func apiGoodRemoveHandler(document http.ResponseWriter, request *http.Request) {

	if request.Method != "POST" {
		log.DebugStr(os.Stdout, "POST以外のGoodRemoveリクエスト")

		writeStruct(document, apiGoodRemoveMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "POSTを使用してください。",
			},
		}, 403)
		return
	}

	programId, err := strconv.Atoi(request.FormValue("p"))

	if err != nil {
		log.Debug(os.Stdout, err)

		writeStruct(document, apiGoodRemoveMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "不正なプログラムIDです。",
			},
		}, 403)
		return
	}

	userId := getSessionUser(request)
	if userId == 0 {
		log.DebugStr(os.Stdout, "匿名のGoodRemoveリクエスト")

		writeStruct(document, apiGoodRemoveMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "削除する権限がありません。",
			},
		}, 403)
		return
	}

	var good models.Good
	err = good.LoadByUserAndProgram(userId, programId)

	if err != nil {
		log.Fatal(os.Stdout, err)

		writeStruct(document, apiGoodRemoveMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "内部エラーが発生しました。",
			},
		}, 500)
		return
	}

	err = good.Remove()

	if err != nil {
		log.Fatal(os.Stdout, err)

		writeStruct(document, apiGoodRemoveMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "内部エラーが発生しました。",
			},
		}, 500)
		return
	}

	writeStruct(document, apiGoodRemoveMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "削除に成功しました。",
		},
	}, 200)

}
