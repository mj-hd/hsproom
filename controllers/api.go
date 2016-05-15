package controllers

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"../bot"
	"../config"
	"../models"
	"../utils"
	"../utils/google"
	"../utils/log"
	"../utils/twitter"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

var twitterClient *twitter.Client
var oauthClient *twitter.OAuthClient
var oauth2Client *google.OAuth2Client
var enabledTwitter bool = true
var enabledOAuth bool = true
var enabledOAuth2 bool = true

func apiInit() {
	var err error

	// 2-legged
	twitterClient, err = twitter.NewClient(config.TwitterKey, config.TwitterSecret)
	if err != nil {
		log.Fatal(err)
		log.FatalStr("TwitterAPIへのアクセスに失敗．Twitter連携機能をオフにします．")
		enabledTwitter = false
	}

	// 3-legged
	oauthClient, err = twitter.NewOAuthClient(config.TwitterKey, config.TwitterSecret)
	if err != nil {
		log.Fatal(err)
		log.FatalStr("TwitterAPIへのアクセスに失敗．Twitter連携機能をオフにします．")
		enabledOAuth = false
	}

	oauth2Client, err = google.NewOAuth2Client(config.GoogleKey, config.GoogleSecret)
	if err != nil {
		log.Fatal(err)
		log.FatalStr("Google+APIへのアクセスに失敗．Google連携機能をオフにします．")
		enabledOAuth2 = false
	}
}
func apiDel() {
	twitterClient.Close()
}

func apiHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	writeStruct(document, apiMember{
		Status:  "success",
		Message: "エラーなし",
	}, http.StatusOK)

	return http.StatusOK, nil
}

type apiProgramGoodMember struct {
	*apiMember
}

func apiProgramGoodHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	// プログラムIDの取得
	rawProgramId := request.FormValue("p")
	programId, err := strconv.Atoi(rawProgramId)

	if err != nil {
		return http.StatusBadRequest, errors.New("プログラムIDが不正です。")
	}

	if !models.ExistsProgram(programId) {
		return http.StatusInternalServerError, errors.New("プログラムの読み込みに失敗しました。")
	}

	user := getSessionUser(request)
	if user == 0 {
		return http.StatusBadRequest, errors.New("ログインが必要です。")
	}

	if !models.CanGoodProgram(user, programId) {
		return http.StatusBadRequest, errors.New("いいね!は一回までです。")
	}

	var good models.Good
	good.UserID = user
	good.ProgramID = programId

	_, err = good.Create()
	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, errors.New("いいね!に失敗しました。")
	}

	writeStruct(document, apiProgramGoodMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "保存に成功しました。",
		},
	}, 200)

	return http.StatusOK, nil
}

type apiProgramUpdateMember struct {
	*apiMember
}

func apiProgramUpdateHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	// 入力値のバリデート
	var rawProgram models.RawProgram
	targetFlags := models.ProgramPublished | models.ProgramID | models.ProgramTitle | models.ProgramThumbnail | models.ProgramDescription | models.ProgramStartax | models.ProgramAttachments | models.ProgramSteps | models.ProgramSourcecode | models.ProgramRuntime | models.ProgramResolution

	rawProgram.ID = request.FormValue("id")
	rawProgram.Title = bluemonday.StrictPolicy().Sanitize(request.FormValue("title"))
	rawProgram.Thumbnail = request.FormValue("thumbnail")
	rawProgram.Description = request.FormValue("description")
	rawProgram.Startax = request.FormValue("startax")
	rawProgram.Attachments = request.FormValue("attachments")
	rawProgram.Steps = request.FormValue("steps")
	rawProgram.Sourcecode = request.FormValue("sourcecode")
	rawProgram.ResolutionW = request.FormValue("resolution_w")
	rawProgram.ResolutionH = request.FormValue("resolution_h")
	rawProgram.Runtime = request.FormValue("runtime")
	rawProgram.Published = request.FormValue("published")

	if rawProgram.Steps == "" {
		targetFlags -= models.ProgramSteps
	}
	if rawProgram.ResolutionW == "" || rawProgram.ResolutionH == "" {
		targetFlags -= models.ProgramResolution
	}
	if rawProgram.Sourcecode == "" {
		targetFlags -= models.ProgramSourcecode
	}
	if rawProgram.Thumbnail == "" {
		targetFlags -= models.ProgramThumbnail
	}
	if rawProgram.Startax == "" {
		targetFlags -= models.ProgramStartax
	}

	err = rawProgram.Validate(targetFlags)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// プログラムへ変換
	program, err := rawProgram.ToProgram(targetFlags)
	if err != nil {
		log.Debug(err)
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	// プログラムの確認
	var prevProg models.Program

	err = prevProg.Load(program.ID)
	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("プログラムIDが不正です。")
	}

	// ユーザのチェック
	if getSessionUser(request) != prevProg.UserID {
		return http.StatusBadRequest, errors.New("プログラムを編集する権限がありません。")
	}

	// TODO: エラー処理
	prevProg.LoadThumbnail()
	prevProg.LoadStartax()
	prevProg.LoadAttachments()

	// 以前のプログラムと合成する
	program.CreatedAt = prevProg.CreatedAt
	program.UserID = prevProg.UserID
	program.Good = prevProg.Good
	program.Play = prevProg.Play

	if (targetFlags & models.ProgramThumbnail) != 0 {
		prevProg.Thumbnail.Data = program.Thumbnail.Data
	}
	program.Thumbnail = prevProg.Thumbnail
	if (targetFlags & models.ProgramStartax) != 0 {
		prevProg.Startax.Data = program.Startax.Data
	}
	program.Startax = prevProg.Startax

	for _, att := range prevProg.Attachments {
		for i := 0; i < len(program.Attachments); i++ {
			if att.Name == program.Attachments[i].Name {
				att.Data = program.Attachments[i].Data
				program.Attachments[i] = att

				goto Found
			}
		}

		// NotFound
		att.Remove()

	Found:
	}

	err = program.Update()
	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, errors.New("保存に失敗しました。")
	}

	writeStruct(document, apiProgramUpdateMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "保存に成功しました。",
		},
	}, http.StatusOK)

	return http.StatusOK, nil
}

type apiProgramCreateMember struct {
	*apiMember
	ID int
}

type apiNameValuePair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func apiProgramCreateHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	// 入力値のバリデート
	var rawProgram models.RawProgram
	targetFlags := models.ProgramPublished | models.ProgramTitle | models.ProgramThumbnail | models.ProgramDescription | models.ProgramStartax | models.ProgramAttachments | models.ProgramSteps | models.ProgramSourcecode | models.ProgramRuntime | models.ProgramResolution

	rawProgram.Title = bluemonday.StrictPolicy().Sanitize(request.FormValue("title"))
	rawProgram.Thumbnail = request.FormValue("thumbnail")
	rawProgram.Description = request.FormValue("description")
	rawProgram.Startax = request.FormValue("startax")
	rawProgram.Attachments = request.FormValue("attachments")
	rawProgram.Steps = request.FormValue("steps")
	rawProgram.Sourcecode = request.FormValue("sourcecode")
	rawProgram.ResolutionW = request.FormValue("resolution_w")
	rawProgram.ResolutionH = request.FormValue("resolution_h")
	rawProgram.Runtime = request.FormValue("runtime")
	rawProgram.Published = request.FormValue("published")

	if rawProgram.Steps == "" {
		targetFlags -= models.ProgramSteps
	}
	if rawProgram.ResolutionW == "" || rawProgram.ResolutionH == "" {
		targetFlags -= models.ProgramResolution
	}
	if rawProgram.Sourcecode == "" {
		targetFlags -= models.ProgramSourcecode
	}

	userId := getSessionUser(request)

	// ログインしていない
	if userId == 0 {
		return http.StatusBadRequest, errors.New("ログインする必要があります。")
	}

	err = rawProgram.Validate(targetFlags)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// プログラムへ変換
	program, err := rawProgram.ToProgram(targetFlags)
	if err != nil {
		log.Debug(err)
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	program.UserID = userId

	id, err := program.Create()
	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, errors.New("保存に失敗しました。もう一度お試しください。")
	}

	writeStruct(document, apiProgramCreateMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "保存に成功しました。",
		},
		ID: id,
	}, http.StatusOK)

	if program.Published {
		bot.UpdateTweet("新しいプログラムが投稿されました! #hsproom\n\n " + program.Title + " by " + program.GetUserName() + " " + config.SiteURL + "/program/view/?p=" + strconv.Itoa(id))
	}

	return http.StatusOK, nil
}

type apiProgramDataListMember struct {
	*apiMember
	Names []string
}

func apiProgramDataListHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	programId, err := strconv.Atoi(request.URL.Query().Get("p"))

	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("プログラムIDが不正です。")
	}

	program := models.NewProgram()
	err = program.Load(programId)

	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("プログラムが存在しません。")
	}

	err = program.LoadAttachments()

	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("添付ファイルの取得に失敗しました。")
	}

	var names []string

	for _, att := range program.Attachments {
		names = append(names, att.Name)
	}

	writeStruct(document, apiProgramDataListMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "添付ファイル一覧の取得に成功しました。",
		},
		Names: names,
	}, http.StatusOK)

	return http.StatusOK, nil
}

// jsonじゃないよ
func apiProgramDataHandler(document http.ResponseWriter, request *http.Request) {

	if request.Method != "GET" {
		log.DebugStr("GET以外のDataリクエスト")
		document.WriteHeader(http.StatusBadRequest)
		return
	}

	document.Header().Set("Content-Type", "application/octet-stream")

	rawProgramId := request.URL.Query().Get("p")
	programId, err := strconv.Atoi(rawProgramId)

	if err != nil {
		log.DebugStr("不正なプログラムID")
		document.WriteHeader(http.StatusBadRequest)
		return
	}

	program := models.NewProgram()
	err = program.Load(programId)

	if err != nil {
		log.Debug(err)
		log.DebugStr("プログラムが見つからない")
		document.WriteHeader(http.StatusNotFound)
		return
	}

	fileName := request.URL.Query().Get("f")
	if fileName == "" {
		log.DebugStr("空のDataリクエスト")
		document.WriteHeader(http.StatusNotFound)
		return
	}

	if fileName == "start.ax" {

		program.LoadStartax()

		document.WriteHeader(http.StatusOK)
		document.Write(program.Startax.Data)

		return
	}

	// ファイルを検索する
	err = program.LoadAttachments()

	if err != nil {
		log.DebugStr("添付ファイルの読み込みに失敗")
		document.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, att := range program.Attachments {
		if att.Name == fileName {
			document.WriteHeader(http.StatusOK)
			document.Write(att.Data)
			return
		}
	}

	// ファイルが見つからなかった
	log.DebugStr("指定されたファイルが見つからなかった")
	document.WriteHeader(http.StatusNotFound)
}

// jsonじゃない
func apiProgramThumbnailHandler(document http.ResponseWriter, request *http.Request) {

	if request.Method != "GET" {
		log.DebugStr("GET以外のThumbnailリクエスト")
		document.WriteHeader(http.StatusBadRequest)
		return
	}

	document.Header().Set("Content-Type", "image/png")

	rawProgramId := request.URL.Query().Get("p")
	programId, err := strconv.Atoi(rawProgramId)

	if err != nil {
		log.Debug(err)
		document.WriteHeader(http.StatusBadRequest)
		return
	}

	program := models.NewProgram()

	err = program.Load(programId)

	if err != nil {
		log.Debug(err)
		document.WriteHeader(http.StatusNotFound)
		return
	}

	err = program.LoadThumbnail()

	if err != nil {
		log.Debug(err)
		document.WriteHeader(http.StatusBadRequest)
		return
	}

	document.Write(program.Thumbnail.Data)
}

type apiMarkdownMember struct {
	*apiMember
	Markdown string
}

func apiMarkdownHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

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
					[]byte(text))))}, http.StatusOK)

	return http.StatusOK, nil
}

type apiTwitterSearchMember struct {
	*apiMember
	Tweets twitter.SearchResponse
}

func apiTwitterSearchHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	if !enabledTwitter {
		log.FatalStr("Twitter連携機能がオフです．")
		return http.StatusInternalServerError, errors.New("Twitter連携機能がオフです．")
	}

	rawProgramId := request.URL.Query().Get("p")
	rawNumber := request.URL.Query().Get("n")
	rawOffset := request.URL.Query().Get("o")

	_, err = strconv.Atoi(rawProgramId)
	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("不正なpの値です。")
	}

	number, err := strconv.Atoi(rawNumber)
	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("不正なnの値です。")
	}

	offset, err := strconv.ParseInt(rawOffset, 10, 64)
	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("不正なoの値です。")
	}

	query := "#hsproom"
	query += " #program" + rawProgramId

	tweets, err := twitterClient.SearchTweets(query, number, offset)
	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, errors.New("検索に失敗しました。")
	}

	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	writeStruct(document, apiTwitterSearchMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "検索に成功しました。",
		},
		Tweets: tweets,
	}, http.StatusOK)

	return http.StatusOK, nil
}

type apiOAuthRequestTokenMember struct {
	*apiMember
	AuthURL string
}

func apiTwitterRequestTokenHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	if !enabledOAuth {
		log.FatalStr("Twitterログイン機能がオフです．")
		return http.StatusInternalServerError, errors.New("Twitterログイン機能がオフです．")
	}

	callbackUrl := request.URL.Query().Get("c")

	if callbackUrl == "" {
		log.DebugStr("callback指定のないRequestTokenリクエスト")
		return http.StatusBadRequest, errors.New("コールバック指定が必要です。")
	}

	url, err := oauthClient.GetAuthURL(config.SiteURL + "/api/twitter/access_token/?c=" + url.QueryEscape(callbackUrl))
	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	writeStruct(document, apiOAuthRequestTokenMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "URLの取得に成功しました。",
		},
		AuthURL: url,
	}, http.StatusOK)

	return http.StatusOK, nil
}

// jsonじゃない
func apiTwitterAccessTokenHandler(document http.ResponseWriter, request *http.Request) {

	if request.Method != "GET" {
		log.DebugStr("GET以外のAccessTokenリクエスト")
		document.WriteHeader(http.StatusBadRequest)
		return
	}

	if !enabledOAuth {
		log.FatalStr("Twitterログイン機能がオフです．")
		document.WriteHeader(http.StatusInternalServerError)
		return
	}

	verifier := request.URL.Query().Get("oauth_verifier")
	token := request.URL.Query().Get("oauth_token")

	if verifier == "" || token == "" {
		log.DebugStr("クエリが空")
		document.WriteHeader(http.StatusBadRequest)
		return
	}

	accessToken, err := oauthClient.GetAccessToken(verifier, token)
	if err != nil {
		log.Fatal(err)
		document.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := oauthClient.CheckUserCredentialsAndGetUser(accessToken)

	if err != nil {
		log.Fatal(err)
		document.WriteHeader(http.StatusInternalServerError)
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
			log.Fatal(err)
			document.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {

		id = oldUser.ID
		dbUser.ID = oldUser.ID
		dbUser.CreatedAt = oldUser.CreatedAt

		err = dbUser.Update()

		if err != nil {
			log.Fatal(err)
			document.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	session, err := getSession(request)
	if err != nil {
		log.Fatal(err)
		document.WriteHeader(http.StatusBadRequest)
		return
	}

	session.Values["User"] = id
	session.Save(request, document)

	callbackUrl := request.URL.Query().Get("c")
	if callbackUrl == "" {
		callbackUrl = config.SiteURL + "/"
	}

	http.Redirect(document, request, callbackUrl, http.StatusFound)
	return
}

func apiGoogleRequestTokenHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	if !enabledOAuth2 {
		log.FatalStr("Googleログイン機能がオフです．")
		return http.StatusInternalServerError, errors.New("Googleログイン機能がオフです．")
	}

	url, err := oauth2Client.GetAuthURL(config.SiteURL + "/api/google/access_token/")
	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	session, err := getSession(request)
	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("クッキーが有効ではありません。")
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
	}, http.StatusOK)

	return http.StatusOK, nil
}

func apiGoogleAccessTokenHandler(document http.ResponseWriter, request *http.Request) {

	if request.Method != "GET" {
		log.DebugStr("GET以外のGoogleAccessTokenリクエスト")
		document.WriteHeader(http.StatusBadRequest)
		return
	}

	if !enabledOAuth2 {
		log.FatalStr("Googleログイン機能がオフです．")
		document.WriteHeader(http.StatusInternalServerError)
		return
	}

	verifier := request.URL.Query().Get("state")
	token := request.URL.Query().Get("code")

	if verifier == "" || token == "" {
		log.DebugStr("クエリが空")
		document.WriteHeader(http.StatusBadRequest)
		return
	}

	accessToken, err := oauth2Client.GetToken(verifier, token)
	if err != nil {
		log.Fatal(err)
		document.WriteHeader(http.StatusInternalServerError)
		return
	}

	userinfo, err := oauth2Client.GetUser(accessToken)

	if err != nil {
		log.Fatal(err)
		document.WriteHeader(http.StatusInternalServerError)
		return
	}

	var dbUser models.User
	dbUser.ScreenName = userinfo.IdString
	dbUser.Name = userinfo.Name
	dbUser.Profile = ""
	dbUser.IconURL = userinfo.Picture
	dbUser.Location = userinfo.Locale
	dbUser.Website = userinfo.Link

	var id int

	var oldUser models.User
	err = oldUser.LoadFromScreenName(userinfo.IdString)

	if err != nil {
		id, err = dbUser.Create()

		if err != nil {
			log.Fatal(err)
			document.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {

		id = oldUser.ID
		dbUser.ID = oldUser.ID

		err = dbUser.Update()

		if err != nil {
			log.Fatal(err)
			document.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	session, err := getSession(request)
	if err != nil {
		log.Fatal(err)
		document.WriteHeader(http.StatusBadRequest)
		return
	}

	session.Values["User"] = id

	callbackUrl := session.Values["Callback"].(string)

	session.Values["Callback"] = ""
	session.Save(request, document)

	if callbackUrl == "" {
		callbackUrl = config.SiteURL + "/"
	}

	http.Redirect(document, request, callbackUrl, http.StatusFound)
	return
}

type apiUserInfoMember struct {
	*apiMember
	*models.User
}

func apiUserInfoHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	rawUserId := request.URL.Query().Get("u")
	userId, err := strconv.Atoi(rawUserId)

	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("uの値が不正です。")
	}

	var user models.User
	err = user.Load(userId)

	if err != nil {
		log.DebugStr("存在しないユーザのリクエスト")
		return http.StatusBadRequest, errors.New("存在しないユーザです。")
	}

	writeStruct(document, user, http.StatusOK)
	return http.StatusOK, nil
}

type apiUserProgramListMember struct {
	*apiMember
	Programs     []models.Program
	ProgramCount int
}

func apiUserProgramsHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	userId, err := strconv.Atoi(request.URL.Query().Get("u"))
	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("uの値が不正です。")
	}

	if !models.ExistsUser(userId) {
		log.Fatal(err)
		return http.StatusBadRequest, errors.New("存在しないユーザです。")
	}

	offset, err := strconv.Atoi(request.URL.Query().Get("o"))
	if err != nil {
		offset = 0
	}

	number, err := strconv.Atoi(request.URL.Query().Get("n"))
	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("nの値が不正です。")
	}

	var user models.User
	err = user.Load(userId)

	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	err = user.LoadPrograms()

	if err != nil {
		log.Debug(err)
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	if offset+number > len(user.Programs) {
		number = len(user.Programs) - offset
	}

	writeStruct(document, apiUserProgramListMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "一覧の取得に成功しました。",
		},
		Programs:     user.Programs[offset : offset+number],
		ProgramCount: len(user.Programs),
	}, http.StatusOK)

	return http.StatusOK, nil
}

type apiUserGoodsMember struct {
	*apiMember
	Programs     []models.Program
	ProgramCount int
}

func apiUserGoodsHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	userId, err := strconv.Atoi(request.URL.Query().Get("u"))

	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("uの値が不正です。")
	}

	if !models.ExistsUser(userId) {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("存在しないユーザです。")
	}

	offset, err := strconv.Atoi(request.URL.Query().Get("o"))

	if err != nil {
		log.DebugStr("不正なオフセット。")
		return http.StatusBadRequest, errors.New("oの値が不正です。")
	}

	number, err := strconv.Atoi(request.URL.Query().Get("n"))

	if err != nil {
		log.DebugStr("不正な制限数。")
		return http.StatusBadRequest, errors.New("nの値が不正です。")
	}

	var goods []models.Good

	programCount, err := models.GetGoodListByUser(&goods, userId, offset, number)

	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	programs := make([]models.Program, number)

	for i, good := range goods {

		if good.ProgramID > 0 {
			err = programs[i].Load(good.ProgramID)

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
	}, http.StatusOK)

	return http.StatusOK, nil
}

type apiProgramGoodCountMember struct {
	*apiMember
	GoodCount int
}

func apiProgramGoodCountHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	programId, err := strconv.Atoi(request.URL.Query().Get("p"))

	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("pの値が不正です。")
	}

	if !models.ExistsProgram(programId) {
		log.DebugStr("存在しないプログラムID。")
		return http.StatusBadRequest, errors.New("存在しないプログラムです。")
	}

	goodCount := models.GetGoodCountByProgram(programId)

	writeStruct(document, apiProgramGoodCountMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "いいね数の取得に成功しました。",
		},
		GoodCount: goodCount,
	}, http.StatusOK)

	return http.StatusOK, nil
}

type apiProgramRemoveMember struct {
	*apiMember
}

func apiProgramRemoveHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	programId, err := strconv.Atoi(request.FormValue("p"))

	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("pの値が不正です。")
	}

	userId := getSessionUser(request)
	if userId == 0 {
		log.DebugStr("匿名のProgramRemoveリクエスト")
		return http.StatusBadRequest, errors.New("削除する権限がありません。")
	}

	var program models.Program
	err = program.Load(programId)

	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	if program.UserID != userId {
		log.DebugStr("権限のないProgramRemoveリクエスト")
		return http.StatusBadRequest, errors.New("削除する権限がありません。")
	}

	err = program.Remove()

	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	writeStruct(document, apiProgramRemoveMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "削除に成功しました。",
		},
	}, http.StatusOK)

	return http.StatusOK, nil
}

type apiGoodRemoveMember struct {
	*apiMember
}

func apiGoodRemoveHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	programId, err := strconv.Atoi(request.FormValue("p"))

	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("pの値が不正です。")
	}

	userId := getSessionUser(request)
	if userId == 0 {
		log.DebugStr("匿名のGoodRemoveリクエスト")
		return http.StatusBadRequest, errors.New("削除する権限がありません。")
	}

	var good models.Good
	err = good.LoadByUserAndProgram(userId, programId)

	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	err = good.Remove()

	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	writeStruct(document, apiGoodRemoveMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "削除に成功しました。",
		},
	}, http.StatusOK)

	return http.StatusOK, nil
}

type apiCommentListMember struct {
	*apiMember
	Comments []models.Comment
	Count    int
	MaxID    int
}

func apiCommentListHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	programId, err := strconv.Atoi(request.URL.Query().Get("p"))

	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("pの値が不正です")
	}

	if !models.ExistsProgram(programId) {
		log.DebugStr("存在しないプログラムID")
		return http.StatusBadRequest, errors.New("存在しないプログラムです。")
	}

	number, err := strconv.Atoi(request.URL.Query().Get("n"))
	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("nの値が不正です。")
	}

	offset, err := strconv.Atoi(request.URL.Query().Get("o"))
	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("oの値が不正です。")
	}

	since, err := strconv.Atoi(request.URL.Query().Get("s"))
	if err != nil {
		since = 0
	}

	containsReply, err := strconv.Atoi(request.URL.Query().Get("c"))
	if err != nil {
		containsReply = 0
	}

	var comments []models.Comment
	if containsReply == 0 {
		comments, err = models.GetComments(programId, number, offset, since)
	} else {
		comments, err = models.GetCommentsAndReplies(programId, number, offset, since)
	}
	if err != nil {
		log.Debug(err)
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	count := 0
	for i, _ := range comments {
		count++
		f := func(comment *models.Comment, me interface{}) {
			comment.LoadReplies()
			comment.LoadUser()

			if len(comment.Replies) <= 0 {
				return
			}

			for j, _ := range comment.Replies {
				(me.(func(*models.Comment, interface{})))(&comment.Replies[j], me)
			}

		}
		f(&comments[i], f)

	}

	maxId, err := models.GetCommentsAndRepliesMaxID(programId)
	if err != nil {
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	writeStruct(document, apiCommentListMember{
		apiMember: &apiMember{
			Status:  "success",
			Message: "取得に成功しました。",
		},
		Comments: comments,
		Count:    count,
		MaxID:    maxId,
	}, http.StatusOK)

	return http.StatusOK, nil
}

func apiCommentPostHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	programId, err := strconv.Atoi(request.FormValue("p"))

	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("プログラムIDが不正です。")
	}

	if !models.ExistsProgram(programId) {
		log.DebugStr("存在しないプログラムID")
		return http.StatusBadRequest, errors.New("存在しないプログラムです。")
	}

	replyTo, err := strconv.Atoi(request.FormValue("r"))
	if err != nil {
		replyTo = -1
	}

	message := request.FormValue("m")
	if message == "" {
		log.DebugStr("空のコメント")
		return http.StatusBadRequest, errors.New("コメントが空です。")
	}

	userId := getSessionUser(request)

	var comment models.Comment
	comment.Message = utils.StandardPolicy().Sanitize(message)
	comment.ProgramID = programId
	comment.UserID = userId

	if len(comment.Message) > 200 || len(comment.Message) == 0 {
		log.DebugStr("コメントの文字数が範囲外")
		return http.StatusBadRequest, errors.New("コメントの文字数が範囲外です。")
	}

	if userId != 0 {
		comment.UserName, err = models.GetUserName(userId)
		if err != nil {
			log.Fatal(err)
			return http.StatusInternalServerError, errors.New("コメントの取得に失敗しました。")
		}

	} else {
		userName := request.FormValue("n")

		if userName == "" {
			userName = "名無し"
		}

		comment.UserName = userName
	}

	comment.ReplyTo = replyTo

	err = comment.Create()
	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	writeStruct(document, apiMember{
		Status:  "success",
		Message: "投稿に成功しました。",
	}, http.StatusOK)

	return http.StatusOK, nil
}

func apiCommentDeleteHandler(document http.ResponseWriter, request *http.Request) (status int, err error) {

	commentId, err := strconv.Atoi(request.FormValue("c"))
	if err != nil {
		log.Debug(err)
		return http.StatusBadRequest, errors.New("cの値が不正です。")
	}

	var comment models.Comment
	err = comment.Load(commentId)
	if err != nil {
		log.Debug(err)
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	userId := getSessionUser(request)

	if userId != comment.UserID || userId == 0 {
		log.DebugStr("権限のないコメント削除リクエスト")
		return http.StatusBadRequest, errors.New("削除する権限がありません。")
	}

	err = comment.Remove()
	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, errors.New("内部エラーが発生しました。")
	}

	writeStruct(document, apiMember{
		Status:  "success",
		Message: "削除に成功しました。",
	}, http.StatusOK)

	return http.StatusOK, nil
}
