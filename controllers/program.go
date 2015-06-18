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

func programHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "program.tmpl"

	var goodPrograms []models.Program
	var recentPrograms []models.Program

	i, err := models.GetProgramRankingForAllTime(&goodPrograms, 0, 4)

	if err != nil {
		log.FatalStr("Error At :" + strconv.Itoa(i))
		log.Fatal(err)

		showError(document, request, "エラーが発生しました、管理人へ報告してください。")

		return
	}

	_, err = models.GetProgramListBy(models.ProgramColCreatedAt, &recentPrograms, true, 0, 4)

	if err != nil {

		log.Fatal(err)

		showError(document, request, "エラーが発生しました、管理人へ知らせてください。")

		return
	}

	err = tmpl.Render(document, programMember{
		DefaultMember: &templates.DefaultMember{
			Title:  config.SiteTitle,
			UserID: getSessionUser(request),
		},
		GoodPrograms:   goodPrograms,
		RecentPrograms: recentPrograms,
	})

	if err != nil {
		log.Fatal(err)

		showError(document, request, "ページの表示に失敗しました。管理人へ報告してください。")

		return
	}
}

type programListMember struct {
	*templates.DefaultMember
	Programs []models.Program
}

func programListHandler(document http.ResponseWriter, request *http.Request) {

	sortKey := request.URL.Query().Get("k")
	order := request.URL.Query().Get("o")

	// レイアウトを指定
	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programList.tmpl"

	// プログラムの一覧を取得
	var programs []models.Program
	var err error
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
		log.Fatal(err)
		showError(document, request, "ページの表示に失敗しました。管理人へ問い合わせてください。")
		return
	}

	err = tmpl.Render(document, programListMember{
		DefaultMember: &templates.DefaultMember{
			Title:  "プログラム一覧",
			UserID: getSessionUser(request),
		},
		Programs: programs,
	})
	if err != nil {
		log.Fatal(err)
		showError(document, request, "ページの表示に失敗しました。管理人へ問い合わせてください。")
		return
	}

}

type programViewMember struct {
	*templates.DefaultMember
	Program         models.Program
	RelatedPrograms []models.Program
}

func programViewHandler(document http.ResponseWriter, request *http.Request) {

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
		showError(document, request, "リクエストが不正です。")

		document.WriteHeader(400)

		return
	}

	// プログラムを取得
	var program models.Program

	err = program.Load(programId)
	if err != nil {

		log.Info(err)
		showError(document, request, "プログラムが存在しません。")

		document.WriteHeader(404)

		return
	}

	userId := getSessionUser(request)

	if program.Published {

		err = models.PlayProgram(program.ID)

		if err != nil {
			log.Fatal(err)

			showError(document, request, "エラーが発生しました。")
			return
		}
	} else {
		if program.UserID != userId {
			log.Debug(errors.New("非公開のプログラムへのアクセス"))

			showError(document, request, "非公開のプログラムです。")
			return
		}
	}

	var related []models.Program
	err = models.GetProgramListRelatedTo(&related, program.Title, 10)

	if err != nil {
		related = make([]models.Program, 0)
	}

	err = tmpl.Render(document, programViewMember{
		DefaultMember: &templates.DefaultMember{
			Title:  program.Title + " - " + config.SiteTitle,
			UserID: getSessionUser(request),
		},
		Program:         program,
		RelatedPrograms: related,
	})
	if err != nil {
		log.Fatal(err)
		showError(document, request, "ページの表示に失敗しました。管理人へ問い合わせてください。")
	}

}

type programPostMember struct {
	*templates.DefaultMember
}

func programPostHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programPost.tmpl"

	user := getSessionUser(request)
	if user == 0 {
		// TODO: ログインを促す

		return
	}

	err := tmpl.Render(document, programPostMember{
		DefaultMember: &templates.DefaultMember{
			Title:  "プログラムの投稿 - " + config.SiteTitle,
			UserID: user,
		},
	})
	if err != nil {
		log.Fatal(err)
		showError(document, request, "ページの表示に失敗しました。管理人へ問い合わせてください。")
	}

}

type programEditMember struct {
	*templates.DefaultMember
	Program *models.Program
}

func programEditHandler(document http.ResponseWriter, request *http.Request) {

	// プログラムIdの取得
	rawProgramId := request.URL.Query().Get("p")
	programId, err := strconv.Atoi(rawProgramId)

	if err != nil {
		log.Debug(err)

		showError(document, request, "プログラムが見つかりません。")

		return
	}

	// テンプレートの設定
	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programEdit.tmpl"

	// ユーザの取得
	user := getSessionUser(request)
	if user == 0 {
		// TODO: ログインを促す

		return
	}

	// プログラムの読み込み
	program := models.NewProgram()
	err = program.Load(programId)

	if err != nil {

		log.Debug(err)

		showError(document, request, "プログラムの読み込みに失敗しました。管理人へ問い合わせてください。")

		return
	}

	if program.UserID != user {
		log.DebugStr("権限のない編集画面へのアクセス")

		showError(document, request, "プログラムの編集権限がありません。")

		return
	}

	if program.Sourcecode != "" {
		from := request.URL.Query().Get("f")

		if from != "source" {
			http.Redirect(document, request, "/source/edit/?p="+strconv.Itoa(program.ID), 303)
			return
		}
	}

	err = program.LoadThumbnail()

	if err != nil {
		log.DebugStr("サムネイル画像の読み込みに失敗しました。")

		showError(document, request, "サムネイル画像の読み込みに失敗しました。")

		return
	}

	err = program.LoadAttachments()

	if err != nil {
		log.DebugStr("添付ファイルの読み込みに失敗しました")

		showError(document, request, "添付ファイルの読み込みに失敗しました")

		return
	}

	// 表示
	err = tmpl.Render(document, programEditMember{
		DefaultMember: &templates.DefaultMember{
			Title:  program.Title + " - " + config.SiteTitle,
			UserID: user,
		},
		Program: program,
	})

	if err != nil {

		log.Debug(err)

		showError(document, request, "ページの読み込みに失敗しました。管理人へ問い合わせてください。")

		return
	}

}

func programCreateHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programCreate.tmpl"

	err := tmpl.Render(document, &templates.DefaultMember{
		Title:  "新規プログラム - " + config.SiteTitle,
		UserID: getSessionUser(request),
	})

	if err != nil {
		log.Fatal(err)

		showError(document, request, "エラーが発生しました。管理人へ報告してください。")
	}

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

func programSearchHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programSearch.tmpl"

	var programs []models.Program
	var sortKey models.ProgramColumn

	queryWord := bluemonday.UGCPolicy().Sanitize(request.URL.Query().Get("q"))
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
		log.Fatal(err)

		showError(document, request, "エラーが発生しました。")

		return
	}

	maxPage := 0
	if i != 0 {
		maxPage = i / 10
		if (maxPage % 10) == 0 {
			maxPage--
		}
	}

	err = tmpl.Render(document, programSearchMember{
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
	if err != nil {
		log.Fatal(err)

		showError(document, request, "ページの表示に失敗しました。管理人へ報告してください。")
		return
	}
}

type programRankingMember struct {
	*templates.DefaultMember
	Programs     []models.Program
	CurPage      int
	MaxPage      int
	ProgramCount int
	Period       string
}

func programRankingDailyHandler(document http.ResponseWriter, request *http.Request) {

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
		log.Fatal(err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	maxPage := i / 10
	if (i % 10) == 0 {
		maxPage--
	}

	err = tmpl.Render(document, programRankingMember{
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

	if err != nil {
		log.Fatal(err)

		showError(document, request, "ページの表示に失敗しました。管理人へ報告してください。")
		return
	}
}

func programRankingMonthlyHandler(document http.ResponseWriter, request *http.Request) {

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
		log.Fatal(err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	maxPage := i / 10
	if (i % 10) == 0 {
		maxPage--
	}

	err = tmpl.Render(document, programRankingMember{
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

	if err != nil {
		log.Fatal(err)

		showError(document, request, "ページの表示に失敗しました。管理人へ報告してください。")
		return
	}

}

func programRankingWeeklyHandler(document http.ResponseWriter, request *http.Request) {

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
		log.Fatal(err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	maxPage := i / 10
	if (i % 10) == 0 {
		maxPage--
	}

	err = tmpl.Render(document, programRankingMember{
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

	if err != nil {
		log.Fatal(err)

		showError(document, request, "ページの表示に失敗しました。管理人へ報告してください。")
		return
	}

}

func programRankingAllTimeHandler(document http.ResponseWriter, request *http.Request) {

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
		log.Fatal(err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	maxPage := i / 10
	if (i % 10) == 0 {
		maxPage--
	}

	err = tmpl.Render(document, programRankingMember{
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

	if err != nil {
		log.Fatal(err)

		showError(document, request, "ページの表示に失敗しました。管理人へ報告してください。")
		return
	}

}
