package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/microcosm-cc/bluemonday"

	"../config"
	"../models"
	"../templates"
	"../utils/log"
)

type programMember struct {
	*templates.DefaultMember
	GoodPrograms   []models.Program
	RecentPrograms []models.Program
}

func programHandler(document http.ResponseWriter, request *http.Request) (err error) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "program.tmpl"

	var goodPrograms []models.Program
	var recentPrograms []models.Program

	_, err = models.GetProgramRankingForAllTime(&goodPrograms, 0, 4)

	if err != nil {
		return errors.New("人気順プログラムの取得に失敗: \r\n" + err.Error())
	}

	_, err = models.GetProgramListBy(models.ProgramColCreatedAt, &recentPrograms, true, 0, 4)

	if err != nil {
		return errors.New("新着順プログラムの取得に失敗: \r\n" + err.Error())
	}

	return tmpl.Render(document, programMember{
		DefaultMember: &templates.DefaultMember{
			Title:  config.SiteTitle,
			UserID: getSessionUser(request),
		},
		GoodPrograms:   goodPrograms,
		RecentPrograms: recentPrograms,
	})
}

type programListMember struct {
	*templates.DefaultMember
	Programs []models.Program
}

func programListHandler(document http.ResponseWriter, request *http.Request) (err error) {

	sortKey := request.URL.Query().Get("k")
	order := request.URL.Query().Get("o")

	// レイアウトを指定
	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programList.tmpl"

	// プログラムの一覧を取得
	var programs []models.Program
	isDesc := (order == "d")

	// sortKeyでふるい分け
	var keyColumn models.ProgramColumn
	switch sortKey {
	case "c": // Created
		keyColumn = models.ProgramColCreatedAt
	case "g": // Good
		keyColumn = models.ProgramColGood
	default:
		keyColumn = models.ProgramColCreatedAt
	}

	_, err = models.GetProgramListBy(keyColumn, &programs, isDesc, 0, 10)

	if err != nil {
		return errors.New("プログラム一覧の取得に失敗: \r\n" + err.Error())
	}

	return tmpl.Render(document, programListMember{
		DefaultMember: &templates.DefaultMember{
			Title:  "プログラム一覧",
			UserID: getSessionUser(request),
		},
		Programs: programs,
	})
}

type programViewMember struct {
	*templates.DefaultMember
	Program         models.Program
	RelatedPrograms []models.Program
}

func programViewHandler(document http.ResponseWriter, request *http.Request) (err error) {

	var tmpl templates.Template

	// スマホ
	if strings.Contains(request.UserAgent(), "iPod") ||
		strings.Contains(request.UserAgent(), "iPhone") ||
		strings.Contains(request.UserAgent(), "Android") {

		tmpl.Layout = "empty.tmpl"
		tmpl.Template = "programViewSP.tmpl"

	} else {

		tmpl.Layout = "default.tmpl"
		tmpl.Template = "programView.tmpl"

	}

	// プログラムIDの取得
	rawProgramId := request.URL.Query().Get("p")
	programId, err := strconv.Atoi(rawProgramId)

	if err != nil {

		log.Info(err)

		showError(document, request, "リクエストが不正です。Request:"+rawProgramId)

		return nil
	}

	// プログラムを取得
	var program models.Program

	err = program.Load(programId)
	if err != nil {

		log.Info(err)

		showError(document, request, "プログラムが存在しません。")

		return nil
	}

	userId := getSessionUser(request)

	if program.Published {

		err = models.PlayProgram(program.ID)

		if err != nil {
			log.FatalStr("プレイ回数を加算できませんでした。ProgramID:" + string(program.ID))
		}
	} else {
		if program.UserID != userId {
			log.DebugStr("非公開のプログラムへのアクセス。ProgramID:" + string(program.ID))

			showError(document, request, "非公開のプログラムです。")
			return nil
		}
	}

	var related []models.Program
	err = models.GetProgramListRelatedTo(&related, program.Title, 10)

	if err != nil {
		related = make([]models.Program, 0)
	}

	return tmpl.Render(document, programViewMember{
		DefaultMember: &templates.DefaultMember{
			Title:  program.Title + " - " + config.SiteTitle,
			UserID: getSessionUser(request),
		},
		Program:         program,
		RelatedPrograms: related,
	})
}

type programRemoteViewMember struct {
	*templates.DefaultMember
	Program models.Program
}

func programRemoteViewHandler(document http.ResponseWriter, request *http.Request) {
	var tmpl templates.Template
	tmpl.Layout = "empty.tmpl"
	tmpl.Template = "programRemoteView.tmpl"

	rawProgramId := request.URL.Query().Get("p")
	programId, err := strconv.Atoi(rawProgramId)

	if err != nil {

		log.Info(err)

		document.WriteHeader(403)
		document.Write([]byte("リクエストが不正です"))

		return
	}

	var program models.Program

	err = program.Load(programId)
	if err != nil {

		log.Info(err)

		document.WriteHeader(404)
		document.Write([]byte("プログラムが見つかりません"))

		return
	}

	if !program.Published {
		document.WriteHeader(503)
		document.Write([]byte("プログラムが非公開です"))
		return
	}

	err = tmpl.Render(document, programViewMember{
		DefaultMember: &templates.DefaultMember{
			Title:  program.Title,
			UserID: getSessionUser(request),
		},
		Program: program,
	})
	if err != nil {
		document.WriteHeader(503)
		document.Write([]byte("サーバでエラーが発生しています"))
	}

	return
}

type programPostMember struct {
	*templates.DefaultMember
}

func programPostHandler(document http.ResponseWriter, request *http.Request) (err error) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programPost.tmpl"

	user := getSessionUser(request)
	if user == 0 {

		http.Redirect(document, request, "/user/login/", http.StatusFound)

		return nil
	}

	return tmpl.Render(document, programPostMember{
		DefaultMember: &templates.DefaultMember{
			Title:  "プログラムの投稿 - " + config.SiteTitle,
			UserID: user,
		},
	})
}

type programEditMember struct {
	*templates.DefaultMember
	Program              *models.Program
	ThumbnailLimitSize   int
	StartaxLimitSize     int
	AttachmentsLimitSize int
}

func programEditHandler(document http.ResponseWriter, request *http.Request) (err error) {

	// プログラムIdの取得
	rawProgramId := request.URL.Query().Get("p")
	programId, err := strconv.Atoi(rawProgramId)

	if err != nil {
		log.Debug(err)

		showError(document, request, "リクエストが不正です。")

		return nil
	}

	// テンプレートの設定
	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programEdit.tmpl"

	// ユーザの取得
	user := getSessionUser(request)
	if user == 0 {

		http.Redirect(document, request, "/user/login/", http.StatusFound)

		return nil
	}

	// プログラムの読み込み
	program := models.NewProgram()
	err = program.Load(programId)

	if err != nil {

		log.Debug(err)

		showError(document, request, "プログラムの読み込みに失敗しました。管理人へ問い合わせてください。")

		return nil
	}

	if program.UserID != user {
		log.DebugStr("権限のない編集画面へのアクセス。ProgramID:" + string(program.ID) + " UserID:" + string(user))

		showError(document, request, "プログラムの編集権限がありません。")

		return nil
	}

	if program.Sourcecode != "" {
		from := request.URL.Query().Get("f")

		if from != "source" {
			http.Redirect(document, request, "/source/edit/?p="+strconv.Itoa(program.ID), http.StatusFound)
			return nil
		}
	}

	err = program.LoadThumbnail()

	if err != nil {
		log.DebugStr("サムネイル画像の読み込みに失敗しました。ProgramID:" + string(program.ID))

		showError(document, request, "サムネイル画像の読み込みに失敗しました。")

		return nil
	}

	err = program.LoadAttachments()

	if err != nil {
		log.DebugStr("添付ファイルの読み込みに失敗しました。ProgramID:" + string(program.ID))

		showError(document, request, "添付ファイルの読み込みに失敗しました")

		return nil
	}

	// 表示
	return tmpl.Render(document, programEditMember{
		DefaultMember: &templates.DefaultMember{
			Title:  program.Title + " - " + config.SiteTitle,
			UserID: user,
		},
		Program:              program,
		ThumbnailLimitSize:   config.ThumbnailLimitSize,
		StartaxLimitSize:     config.StartaxLimitSize,
		AttachmentsLimitSize: config.AttachmentsLimitSize,
	})

}

type programCreateMember struct {
	*templates.DefaultMember
	ThumbnailLimitSize   int
	StartaxLimitSize     int
	AttachmentsLimitSize int
}

func programCreateHandler(document http.ResponseWriter, request *http.Request) (err error) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programCreate.tmpl"

	return tmpl.Render(document, programCreateMember{
		DefaultMember: &templates.DefaultMember{
			Title:  "新規プログラム - " + config.SiteTitle,
			UserID: getSessionUser(request),
		},
		ThumbnailLimitSize:   config.ThumbnailLimitSize,
		StartaxLimitSize:     config.StartaxLimitSize,
		AttachmentsLimitSize: config.AttachmentsLimitSize,
	})
}

type programSearchMember struct {
	*templates.DefaultMember
	Query        string
	Sort         string
	CurPage      int
	MaxPage      int
	Programs     []models.Program
	ProgramCount int
}

func programSearchHandler(document http.ResponseWriter, request *http.Request) (err error) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programSearch.tmpl"

	var programs []models.Program
	var sortKey models.ProgramColumn

	queryWord := bluemonday.StrictPolicy().Sanitize(request.URL.Query().Get("q"))
	rawSortKey := request.URL.Query().Get("s")
	page, err := strconv.Atoi(request.URL.Query().Get("p"))

	if err != nil {
		page = 0
	}
	if rawSortKey == "" {
		rawSortKey = "c"
	}

	switch rawSortKey {
	case "c":
		sortKey = models.ProgramColCreatedAt
	case "m":
		sortKey = models.ProgramColUpdatedAt
	case "g":
		sortKey = models.ProgramColGood
	default:
		sortKey = models.ProgramColCreatedAt
	}

	i, err := models.GetProgramListByQuery(&programs, queryWord, sortKey, true, 10, page*10)

	if err != nil {
		return errors.New("プログラム一覧の取得に失敗:\r\n" + err.Error())
	}

	maxPage := 0
	if i != 0 {
		maxPage = i / 10
		if (maxPage % 10) == 0 {
			maxPage--
		}
	}

	return tmpl.Render(document, programSearchMember{
		DefaultMember: &templates.DefaultMember{
			Title:  "プログラムの検索 - " + config.SiteTitle,
			UserID: getSessionUser(request),
		},
		Query:        queryWord,
		Sort:         rawSortKey,
		CurPage:      page,
		MaxPage:      maxPage,
		Programs:     programs,
		ProgramCount: i,
	})
}

type programRankingMember struct {
	*templates.DefaultMember
	Programs     []models.Program
	CurPage      int
	MaxPage      int
	ProgramCount int
	Period       string
}

func programRankingDailyHandler(document http.ResponseWriter, request *http.Request) (err error) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programRanking.tmpl"

	page, err := strconv.Atoi(request.URL.Query().Get("p"))

	if err != nil {
		page = 0
	}

	var programs []models.Program

	i, err := models.GetProgramRankingForDay(&programs, page*10, 10)

	if err != nil {
		return errors.New("日間ランキングの取得に失敗: \r\n" + err.Error())
	}

	maxPage := i / 10
	if (i % 10) == 0 {
		maxPage--
	}

	return tmpl.Render(document, programRankingMember{
		DefaultMember: &templates.DefaultMember{
			Title:  "日間ランキング - " + config.SiteTitle,
			UserID: getSessionUser(request),
		},
		Programs:     programs,
		CurPage:      page,
		MaxPage:      maxPage,
		ProgramCount: i,
		Period:       "daily",
	})
}

func programRankingMonthlyHandler(document http.ResponseWriter, request *http.Request) (err error) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programRanking.tmpl"

	page, err := strconv.Atoi(request.URL.Query().Get("p"))

	if err != nil {
		page = 0
	}

	var programs []models.Program

	i, err := models.GetProgramRankingForMonth(&programs, page*10, 10)

	if err != nil {
		return errors.New("月間ランキングの取得に失敗: \r\n" + err.Error())
	}

	maxPage := i / 10
	if (i % 10) == 0 {
		maxPage--
	}

	return tmpl.Render(document, programRankingMember{
		DefaultMember: &templates.DefaultMember{
			Title:  "月間ランキング - " + config.SiteTitle,
			UserID: getSessionUser(request),
		},
		Programs:     programs,
		CurPage:      page,
		MaxPage:      maxPage,
		ProgramCount: i,
		Period:       "monthly",
	})
}

func programRankingWeeklyHandler(document http.ResponseWriter, request *http.Request) (err error) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programRanking.tmpl"

	page, err := strconv.Atoi(request.URL.Query().Get("p"))

	if err != nil {
		page = 0
	}

	var programs []models.Program

	i, err := models.GetProgramRankingForWeek(&programs, page*10, 10)

	if err != nil {
		return errors.New("週間ランキングの取得に失敗: \r\n" + err.Error())
	}

	maxPage := i / 10
	if (i % 10) == 0 {
		maxPage--
	}

	return tmpl.Render(document, programRankingMember{
		DefaultMember: &templates.DefaultMember{
			Title:  "週間ランキング - " + config.SiteTitle,
			UserID: getSessionUser(request),
		},
		Programs:     programs,
		CurPage:      page,
		MaxPage:      maxPage,
		ProgramCount: i,
		Period:       "weekly",
	})
}

func programRankingAllTimeHandler(document http.ResponseWriter, request *http.Request) (err error) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programRanking.tmpl"

	page, err := strconv.Atoi(request.URL.Query().Get("p"))

	if err != nil {
		page = 0
	}

	var programs []models.Program

	i, err := models.GetProgramRankingForAllTime(&programs, page*10, 10)

	if err != nil {
		return errors.New("総合ランキングの取得に失敗: \r\n" + err.Error())
	}

	maxPage := i / 10
	if (i % 10) == 0 {
		maxPage--
	}

	return tmpl.Render(document, programRankingMember{
		DefaultMember: &templates.DefaultMember{
			Title:  "総合ランキング" + config.SiteTitle,
			UserID: getSessionUser(request),
		},
		Programs:     programs,
		CurPage:      page,
		MaxPage:      maxPage,
		ProgramCount: i,
		Period:       "alltime",
	})
}
