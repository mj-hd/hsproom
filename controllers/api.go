package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"hsproom/config"
	"hsproom/models"
	"hsproom/utils"
	"hsproom/utils/twitter"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

var twitterClient *twitter.Client
var oauthClient *twitter.OAuthClient

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
}
func apiDel() {
	twitterClient.Close()
}

func writeStruct(document http.ResponseWriter, s interface{}, httpStatus int) {

	var err error

	document.Header().Set("Content-Type", "application/json")
	jso, err := json.Marshal(s)

	if err != nil {

		utils.PromulgateFatal(os.Stdout, err)

		document.WriteHeader(500)
		document.Write([]byte("{ \"Status\" : \"error\", \"Message\" : \"なんか変…\" }"))

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

		utils.PromulgateDebugStr(os.Stdout, "POST以外のGoodリクエスト")

		writeStruct(document, apiProgramGoodMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "POST以外のメソッドです。",
			},
		}, 400)

		return
	}

	// プログラムIDの取得
	rawProgramId := request.FormValue("pid")
	programId, err := strconv.Atoi(rawProgramId)

	if err != nil {

		utils.PromulgateDebugStr(os.Stdout, "不正なプログラムID "+string(programId))
		utils.PromulgateDebug(os.Stdout, err)

		writeStruct(document, apiProgramGoodMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "プログラムIDが不正です。",
			},
		}, 400)

		return
	}

	var program models.ProgramInfo
	err = program.Load(programId)

	if err != nil {

		utils.PromulgateFatalStr(os.Stdout, "プログラム["+string(programId)+"]の読み込みに失敗")
		utils.PromulgateFatal(os.Stdout, err)

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

		utils.PromulgateDebugStr(os.Stdout, "匿名のGoodリクエスト")

		writeStruct(document, apiProgramGoodMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "ログインが必要です。",
			},
		}, 400)

		return
	}

	err = program.GiveGood()
	if err != nil {

		utils.PromulgateFatal(os.Stdout, err)
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

		utils.PromulgateDebugStr(os.Stdout, "POST以外のUpdateリクエスト")

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
	targetFlags := models.ProgramId | models.ProgramTitle | models.ProgramThumbnail | models.ProgramDescription | models.ProgramStartax | models.ProgramSize | models.ProgramAttachments

	rawProgram.Id = request.FormValue("id")
	rawProgram.Title = request.FormValue("title")
	rawProgram.Thumbnail = request.FormValue("thumbnail")
	rawProgram.Description = request.FormValue("description")
	rawProgram.Startax = request.FormValue("startax")
	rawProgram.Size = request.FormValue("size")
	rawProgram.Attachments = request.FormValue("attachments")

	err := rawProgram.Validate(targetFlags)
	if err != nil {
		utils.PromulgateDebug(os.Stdout, err)

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
		utils.PromulgateDebug(os.Stdout, err)

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
		utils.PromulgateDebug(os.Stdout, err)

		writeStruct(document, apiProgramUpdateMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "プログラムIDが不正です。",
			},
		}, 400)

		return
	}

	// ユーザのチェック
	if getSessionUser(request) != prevProgInfo.UserId {
		utils.PromulgateDebugStr(os.Stdout, "プログラムの権限のない変更")

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
	prevProgInfo.Thumbnail = program.Thumbnail
	prevProgInfo.Description = program.Description

	// 以前のプログラムと合成する
	program.ProgramInfo = &prevProgInfo

	err = program.Update()
	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)

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
}

type apiNameValuePair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func apiProgramCreateHandler(document http.ResponseWriter, request *http.Request) {

	// メソッドの確認
	if request.Method != "POST" {

		utils.PromulgateDebugStr(os.Stdout, "POST以外のCreateリクエスト")

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
	targetFlags := models.ProgramTitle | models.ProgramUser | models.ProgramThumbnail | models.ProgramDescription | models.ProgramStartax | models.ProgramAttachments

	rawProgram.Title = request.FormValue("title")
	rawProgram.Thumbnail = request.FormValue("thumbnail")
	rawProgram.Description = request.FormValue("description")
	rawProgram.Startax = request.FormValue("startax")
	rawProgram.Attachments = request.FormValue("attachments")

	userId := getSessionUser(request)

	// ログインしていない
	if userId == 0 {
		utils.PromulgateDebugStr(os.Stdout, "匿名のCreateリクエスト")

		writeStruct(document, apiProgramCreateMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "ログインする必要があります。",
			},
		}, 400)

		return
	}

	var userName string
	userName, err := models.GetUserName(userId)

	rawProgram.User = userName

	if err != nil {
		utils.PromulgateDebug(os.Stdout, err)

		writeStruct(document, apiProgramCreateMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: err.Error(),
			},
		}, 400)

		return
	}

	err = rawProgram.Validate(targetFlags)
	if err != nil {
		utils.PromulgateDebug(os.Stdout, err)

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
		utils.PromulgateDebug(os.Stdout, err)

		writeStruct(document, apiProgramCreateMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "内部エラーが発生しました。",
			},
		}, 500)

		return
	}

	program.UserId = userId

	program.Size = len(program.Startax)

	_, err = program.Create()
	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)

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
	}, 200)

}

type apiProgramDataListMember struct {
	*apiMember
	Names []string
}

func apiProgramDataListHandler(document http.ResponseWriter, request *http.Request) {

	if request.Method != "GET" {

		utils.PromulgateDebugStr(os.Stdout, "GET以外のDataListリクエスト")

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

		utils.PromulgateDebugStr(os.Stdout, "プログラムIDが不正")

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

		utils.PromulgateDebug(os.Stdout, err)

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

		utils.PromulgateDebugStr(os.Stdout, "GET以外のDataリクエスト")

		document.WriteHeader(400)

		return
	}

	document.Header().Set("Content-Type", "application/octet-stream")

	rawProgramId := request.URL.Query().Get("pid")
	programId, err := strconv.Atoi(rawProgramId)

	if err != nil {

		utils.PromulgateDebugStr(os.Stdout, "不正なプログラムID")

		document.WriteHeader(400)

		return
	}

	program := models.NewProgram()
	err = program.Load(programId)

	if err != nil {

		utils.PromulgateDebug(os.Stdout, err)
		utils.PromulgateDebugStr(os.Stdout, "プログラムが見つからない")

		document.WriteHeader(404)

		return
	}

	fileName := request.URL.Query().Get("f")
	if fileName == "" {

		utils.PromulgateDebugStr(os.Stdout, "空のDataリクエスト")

		document.WriteHeader(404)

		return
	}

	if fileName == "startax" {

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
	utils.PromulgateDebugStr(os.Stdout, "指定されたファイルが見つからなかった")
	document.WriteHeader(404)
}

// jsonじゃない
func apiProgramThumbnailHandler(document http.ResponseWriter, request *http.Request) {

	if request.Method != "GET" {

		utils.PromulgateDebugStr(os.Stdout, "GET以外のThumbnailリクエスト")

		document.WriteHeader(400)

		return
	}

	document.Header().Set("Content-Type", "image/png")

	rawProgramId := request.URL.Query().Get("p")
	programId, err := strconv.Atoi(rawProgramId)

	if err != nil {

		utils.PromulgateDebug(os.Stdout, err)

		document.WriteHeader(400)

		return
	}

	var programInfo models.ProgramInfo

	err = programInfo.Load(programId)

	if err != nil {

		utils.PromulgateDebug(os.Stdout, err)

		document.WriteHeader(404)

		return
	}

	document.Write(programInfo.Thumbnail)
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

	programName := request.URL.Query().Get("pn")
	rawNumber := request.URL.Query().Get("n")
	rawOffset := request.URL.Query().Get("o")

	number, err := strconv.Atoi(rawNumber)
	if err != nil {
		utils.PromulgateDebug(os.Stdout, err)

		writeStruct(document, apiTwitterSearchMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "不正な数値です。",
			},
		}, 403)
		return
	}

	offset, err := strconv.Atoi(rawOffset)
	if err != nil {
		utils.PromulgateDebug(os.Stdout, err)

		writeStruct(document, apiTwitterSearchMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "不正な数値です。",
			},
		}, 403)
		return
	}

	query := "#hsproom"

	if programName != "" {
		query += " #" + programName
	}

	tweets, err := twitterClient.SearchTweets(query, number, offset)
	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)

		writeStruct(document, apiTwitterSearchMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "検索に失敗しました。",
			},
		}, 500)
		return
	}

	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)

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

type apiTwitterRequestTokenMember struct {
	*apiMember
	AuthURL string
}

func apiTwitterRequestTokenHandler(document http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		utils.PromulgateDebugStr(os.Stdout, "GET以外のRequestTokenリクエスト")

		writeStruct(document, apiTwitterRequestTokenMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "GETを使用してください。",
			},
		}, 403)
		return
	}

	callbackUrl := request.URL.Query().Get("c")

	if callbackUrl == "" {
		utils.PromulgateDebugStr(os.Stdout, "callback指定のないRequestTokenリクエスト")

		writeStruct(document, apiTwitterRequestTokenMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "コールバック指定が必要です。",
			},
		}, 400)
		return
	}

	url, err := oauthClient.GetAuthURL(config.SiteURL + "/api/twitter/access_token/?c=" + url.QueryEscape(callbackUrl))
	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)

		writeStruct(document, apiTwitterRequestTokenMember{
			apiMember: &apiMember{
				Status:  "error",
				Message: "内部エラーが発生しました。",
			},
		}, 500)
		return
	}

	writeStruct(document, apiTwitterRequestTokenMember{
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
		utils.PromulgateDebugStr(os.Stdout, "GET以外のAccessTokenリクエスト")

		document.WriteHeader(403)

		return
	}

	verifier := request.URL.Query().Get("oauth_verifier")
	token := request.URL.Query().Get("oauth_token")

	if verifier == "" || token == "" {
		utils.PromulgateDebugStr(os.Stdout, "クエリが空")

		document.WriteHeader(403)

		return
	}

	accessToken, err := oauthClient.GetAccessToken(verifier, token)
	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)

		document.WriteHeader(500)

		return
	}

	user, err := oauthClient.CheckUserCredentialsAndGetUser(accessToken)

	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)

		document.WriteHeader(500)

		return
	}

	var dbUser models.User
	dbUser.Name = user.Name
	dbUser.Token = accessToken.Token
	dbUser.Secret = accessToken.Secret
	dbUser.Profile = user.Description
	dbUser.IconURL = user.ProfileImageURL
	dbUser.Website = user.URL
	dbUser.Location = user.Location

	var id int

	var oldUser models.User
	err = oldUser.LoadFromName(user.Name)

	if err != nil {
		id, err = dbUser.Create()
		if err != nil {
			utils.PromulgateFatal(os.Stdout, err)

			document.WriteHeader(500)

			return
		}
	} else {

		id = oldUser.Id
		dbUser.Id = oldUser.Id

		err = dbUser.Update()

		if err != nil {
			utils.PromulgateFatal(os.Stdout, err)

			document.WriteHeader(500)

			return
		}

	}

	session, err := sessionStore.Get(request, "go-wiki")
	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)

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

type apiUserInfoMember struct {
	*apiMember
	*models.User
}

func apiUserInfoHandler(document http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		utils.PromulgateDebugStr(os.Stdout, "GET以外のUserInfoリクエスト")

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
		utils.PromulgateDebugStr(os.Stdout, "ユーザ指定のないリクエスト")

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
		utils.PromulgateDebugStr(os.Stdout, "存在しないユーザのリクエスト")

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
