package controllers

import (
	"net/http"
	"os"
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
	GoodPrograms   []models.ProgramInfo
	RecentPrograms []models.ProgramInfo
}

func programHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "program.tmpl"

	var goodPrograms []models.ProgramInfo
	var recentPrograms []models.ProgramInfo

	i, err := models.GetProgramRankingForAllTime(&goodPrograms, 0, 4)

	if err != nil {
		log.FatalStr(os.Stdout, "Error At :"+strconv.Itoa(i))
		log.Fatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました、管理人へ報告してください。")

		return
	}

	_, err = models.GetProgramListBy(models.ProgramColCreated, &recentPrograms, true, 0, 4)

	if err != nil {

		log.Fatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました、管理人へ知らせてください。")

		return
	}

	err = tmpl.Render(document, programMember{
		DefaultMember: &templates.DefaultMember{
			Title: config.SiteTitle,
			User:  getSessionUser(request),
		},
		GoodPrograms:   goodPrograms,
		RecentPrograms: recentPrograms,
	})

	if err != nil {
		log.Fatal(os.Stdout, err)

		showError(document, request, "ページの表示に失敗しました。管理人へ報告してください。")

		return
	}
}

type programListMember struct {
	*templates.DefaultMember
	Programs []models.ProgramInfo
}

func programListHandler(document http.ResponseWriter, request *http.Request) {

	sortKey := request.URL.Query().Get("k")
	order := request.URL.Query().Get("o")

	// レイアウトを指定
	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programList.tmpl"

	// プログラムの一覧を取得
	var programs []models.ProgramInfo
	var err error
	isDesc := (order == "d")

	// sortKeyでふるい分け
	var keyColumn models.ProgramColumn
	switch sortKey {
	case "c": // Created
		keyColumn = models.ProgramColCreated
	case "g": // Good
		keyColumn = models.ProgramColGood
	default:
		keyColumn = models.ProgramColCreated
	}

	_, err = models.GetProgramListBy(keyColumn, &programs, isDesc, 0, 10)

	if err != nil {
		log.Fatal(os.Stdout, err)
		showError(document, request, "ページの表示に失敗しました。管理人へ問い合わせてください。")
		return
	}

	err = tmpl.Render(document, programListMember{
		DefaultMember: &templates.DefaultMember{
			Title: "プログラム一覧",
			User:  getSessionUser(request),
		},
		Programs: programs,
	})
	if err != nil {
		log.Fatal(os.Stdout, err)
		showError(document, request, "ページの表示に失敗しました。管理人へ問い合わせてください。")
		return
	}

}

type programViewMember struct {
	*templates.DefaultMember
	ProgramInfo     models.ProgramInfo
	RelatedPrograms []models.ProgramInfo
}

func programViewHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template

	// スマホ
	if	strings.Contains(request.UserAgent(), "iPod") ||
		strings.Contains(request.UserAgent(), "iPhone") ||
		strings.Contains(request.UserAgent(), "Android") {

		tmpl.Layout   = "empty.tmpl"
		tmpl.Template = "programViewSP.tmpl"
		
	} else {

		tmpl.Layout = "default.tmpl"
		tmpl.Template = "programView.tmpl"

	}

	// プログラムIDの取得
	rawProgramId := request.URL.Query().Get("p")
	programId, err := strconv.Atoi(rawProgramId)

	if err != nil {

		log.Info(os.Stdout, err)
		showError(document, request, "リクエストが不正です。")

		document.WriteHeader(400)

		return
	}

	// プログラムを取得
	var program models.ProgramInfo

	err = program.Load(programId)
	if err != nil {

		log.Info(os.Stdout, err)
		showError(document, request, "プログラムが存在しません。")

		document.WriteHeader(404)

		return
	}

	err = models.PlayProgram(program.Id)

	if err != nil {
		log.Fatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	var related []models.ProgramInfo
	err = models.GetProgramListRelatedTo(&related, program.Title, 10)

	if err != nil {
		related = make([]models.ProgramInfo, 0)
	}

	err = tmpl.Render(document, programViewMember{
		DefaultMember: &templates.DefaultMember{
			Title: program.Title + " - " + config.SiteTitle,
			User:  getSessionUser(request),
		},
		ProgramInfo:     program,
		RelatedPrograms: related,
	})
	if err != nil {
		log.Fatal(os.Stdout, err)
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
			Title: "プログラムの投稿 - " + config.SiteTitle,
			User:  user,
		},
	})
	if err != nil {
		log.Fatal(os.Stdout, err)
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
		log.Debug(os.Stdout, err)

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

		log.Debug(os.Stdout, err)

		showError(document, request, "プログラムの読み込みに失敗しました。管理人へ問い合わせてください。")

		return
	}

	if program.User != user {
		log.DebugStr(os.Stdout, "権限のない編集画面へのアクセス")

		showError(document, request, "プログラムの編集権限がありません。")

		return
	}

	if program.Sourcecode != "" {
		from := request.URL.Query().Get("f")

		if (from != "source") {
			http.Redirect(document, request, "/source/edit/?p="+strconv.Itoa(program.Id), 303)
			return
		}
	}

	// 表示
	err = tmpl.Render(document, programEditMember{
		DefaultMember: &templates.DefaultMember{
			Title: program.Title + " - " + config.SiteTitle,
			User:  user,
		},
		Program: program,
	})

	if err != nil {

		log.Debug(os.Stdout, err)

		showError(document, request, "ページの読み込みに失敗しました。管理人へ問い合わせてください。")

		return
	}

}

func programCreateHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programCreate.tmpl"

	err := tmpl.Render(document, &templates.DefaultMember{
		Title: "新規プログラム - " + config.SiteTitle,
		User:  getSessionUser(request),
	})

	if err != nil {
		log.Fatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました。管理人へ報告してください。")
	}

}

type programSearchMember struct {
	*templates.DefaultMember
	Query        string
	Sort         string
	CurPage      int
	MaxPage      int
	Programs     []models.ProgramInfo
	ProgramCount int
}

func programSearchHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "programSearch.tmpl"

	var programs []models.ProgramInfo
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
		sortKey = models.ProgramColCreated
	case "m":
		sortKey = models.ProgramColModified
	case "g":
		sortKey = models.ProgramColGood
	default:
		sortKey = models.ProgramColCreated
	}

	i, err := models.GetProgramListByQuery(&programs, queryWord, sortKey, true, 10, page*10)

	if err != nil {
		log.Fatal(os.Stdout, err)

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
			Title: "プログラムの検索 - " + config.SiteTitle,
			User:  getSessionUser(request),
		},
		Query:        queryWord,
		Sort:         rawSortKey,
		CurPage:      page,
		MaxPage:      maxPage,
		Programs:     programs,
		ProgramCount: i,
	})
	if err != nil {
		log.Fatal(os.Stdout, err)

		showError(document, request, "ページの表示に失敗しました。管理人へ報告してください。")
		return
	}
}

type programRankingMember struct {
	*templates.DefaultMember
	Programs     []models.ProgramInfo
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

	var programs []models.ProgramInfo

	i, err := models.GetProgramRankingForDay(&programs, page*10, 10)

	if err != nil {
		log.Fatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	maxPage := i / 10
	if (i % 10) == 0 {
		maxPage--
	}

	err = tmpl.Render(document, programRankingMember{
		DefaultMember: &templates.DefaultMember{
			Title: "日間ランキング - " + config.SiteTitle,
			User:  getSessionUser(request),
		},
		Programs:     programs,
		CurPage:      page,
		MaxPage:      maxPage,
		ProgramCount: i,
		Period:       "daily",
	})

	if err != nil {
		log.Fatal(os.Stdout, err)

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

	var programs []models.ProgramInfo

	i, err := models.GetProgramRankingForMonth(&programs, page*10, 10)

	if err != nil {
		log.Fatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	maxPage := i / 10
	if (i % 10) == 0 {
		maxPage--
	}

	err = tmpl.Render(document, programRankingMember{
		DefaultMember: &templates.DefaultMember{
			Title: "月間ランキング - " + config.SiteTitle,
			User:  getSessionUser(request),
		},
		Programs:     programs,
		CurPage:      page,
		MaxPage:      maxPage,
		ProgramCount: i,
		Period:       "monthly",
	})

	if err != nil {
		log.Fatal(os.Stdout, err)

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

	var programs []models.ProgramInfo

	i, err := models.GetProgramRankingForWeek(&programs, page*10, 10)

	if err != nil {
		log.Fatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	maxPage := i / 10
	if (i % 10) == 0 {
		maxPage--
	}

	err = tmpl.Render(document, programRankingMember{
		DefaultMember: &templates.DefaultMember{
			Title: "週間ランキング - " + config.SiteTitle,
			User:  getSessionUser(request),
		},
		Programs:     programs,
		CurPage:      page,
		MaxPage:      maxPage,
		ProgramCount: i,
		Period:       "weekly",
	})

	if err != nil {
		log.Fatal(os.Stdout, err)

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

	var programs []models.ProgramInfo

	i, err := models.GetProgramRankingForAllTime(&programs, page*10, 10)

	if err != nil {
		log.Fatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	maxPage := i / 10
	if (i % 10) == 0 {
		maxPage--
	}

	err = tmpl.Render(document, programRankingMember{
		DefaultMember: &templates.DefaultMember{
			Title: "総合ランキング" + config.SiteTitle,
			User:  getSessionUser(request),
		},
		Programs:     programs,
		CurPage:      page,
		MaxPage:      maxPage,
		ProgramCount: i,
		Period:       "alltime",
	})

	if err != nil {
		log.Fatal(os.Stdout, err)

		showError(document, request, "ページの表示に失敗しました。管理人へ報告してください。")
		return
	}

}
