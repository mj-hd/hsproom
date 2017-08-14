package models

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/lestrrat/go-ngram"

	"../config"
)

func initPrograms() {
	DB.AutoMigrate(&Program{})
	DB.AutoMigrate(&Attachment{})
	DB.AutoMigrate(&Thumbnail{})
	DB.AutoMigrate(&Startax{})
}

type Program struct {
	ID          int
	CreatedAt   time.Time  ``
	UpdatedAt   time.Time  ``
	DeletedAt   *time.Time ``
	Title       string     `sql:"size:100;not null"`
	UserID      int        `sql:"not null;index"`
	UserName    string     `sql:"-"`
	Good        int        `sql:"default:0"`
	Play        int        `sql:"default:0"`
	Description string     `sql:"size:500"`
	Steps       int        `sql:"default:5000"`
	Runtime     string     `sql:"size:10;default:'HSP3Dish'"`
	RuntimeVersion string  `sql:"size:20;default:'hsp3.5b2mod'"`
	Published   bool       `sql:"not null"`
	ResolutionW int        `sql:"default:640"`
	ResolutionH int        `sql:"default:480"`

	Startax     Startax      ``
	Attachments []Attachment ``
	Thumbnail   Thumbnail    ``
	Sourcecode  string       `sql:"type:text"`
	Goods       []Good
}

type Thumbnail struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	ProgramID int `sql:"index"`

	Data []byte `sql:"type:longblob"`
}

type Startax struct {
	ID         int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	ProgramID  int        `sql:"index"`
	Attachment Attachment `gorm:"polymorphic:Owner;"`

	Data []byte `sql:"type:longblob"`
}

type Attachment struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	ProgramID int `sql:"index"`

	Name string `sql:"size:100;not null"`
	Data []byte `sql:"type:longblob;not null"`
}

func (this *Attachment) ToBase64() string {
	return base64.StdEncoding.EncodeToString(this.Data)
}

func (this *Attachment) Load(id int) error {
	return DB.First(this, id).Error
}

func (this *Attachment) Update() error {
	return DB.Save(this).Error
}

func (this *Attachment) Create() error {
	return DB.Create(this).Error
}

func (this *Attachment) Remove() error {
	return DB.Delete(this).Error
}

func NewProgram() *Program {
	return &Program{}
}

func Published(db *gorm.DB) *gorm.DB {
	return db.Where("published = ?", 1)
}

func (this *Program) AfterFind() (err error) {
	this.LoadUserName()
	this.CreatedAt = this.CreatedAt.In(config.JST())
	this.UpdatedAt = this.UpdatedAt.In(config.JST())

	if this.DeletedAt != nil {
		*this.DeletedAt = this.DeletedAt.In(config.JST())
	}

	return nil
}

func (this *Program) Load(id int) error {

	err := DB.First(this, id).Error

	if err != nil {
		return err
	}

	this.LoadUserName()

	return nil
}

func (this *Program) LoadAttachments() error {
	return DB.Model(this).Related(&this.Attachments).Error
}

func (this *Program) LoadThumbnail() error {
	return DB.Model(this).Related(&this.Thumbnail).Error
}

func (this *Program) LoadStartax() error {
	return DB.Model(this).Related(&this.Startax).Error
}

func (this *Program) LoadUserName() {
	this.UserName = this.GetUserName()
}

func (this *Program) Update() error {

	err := DB.Save(this).Error

	return err
}

func (this *Program) Create() (int, error) {

	err := DB.Create(this).Error

	return this.ID, err
}

func (this *Program) Remove() error {
	var err error

	tx := DB.Begin()

	err = tx.Where("program_id = ?", this.ID).Delete(&this.Attachments).Error
	if err != nil && (err != gorm.ErrRecordNotFound) {
		tx.Rollback()
		return err
	}

	err = tx.Model(this).Related(&this.Thumbnail).Delete(&this.Thumbnail).Error
	if err != nil && (err != gorm.ErrRecordNotFound) {
		tx.Rollback()
		return err
	}

	err = tx.Model(this).Related(&this.Startax).Delete(&this.Startax).Error
	if err != nil && (err != gorm.ErrRecordNotFound) {
		tx.Rollback()
		return err
	}

	err = tx.Where("program_id = ?", this.ID).Delete(Comment{}).Error
	if err != nil && (err != gorm.ErrRecordNotFound) {
		tx.Rollback()
		return err
	}

	err = tx.Where("program_id = ?", this.ID).Delete(&this.Goods).Error
	if err != nil && (err != gorm.ErrRecordNotFound) {
		tx.Rollback()
		return err
	}

	err = tx.Delete(this).Error
	if err != nil {
		tx.Rollback()
	}

	return tx.Commit().Error
}

func (this *Program) FindAttachment(name string) (*Attachment, error) {
	for i, att := range this.Attachments {
		if att.Name == name {
			return &this.Attachments[i], nil
		}
	}

	return nil, errors.New("ファイル" + name + "が見つかりませんでした。")
}

func (this *Program) GetUser() (*User, error) {
	var result User
	err := DB.Model(this).Related(&result).Error
	return &result, err
}

func (this Program) GetScreenName() string {

	name, _ := GetUserScreenName(this.UserID)

	return name
}

func (this Program) GetUserName() string {

	name, _ := GetUserName(this.UserID)

	return name
}

type RawProgram struct {
	ID          string
	Title       string
	UserID      string
	Thumbnail   string
	Description string
	Startax     string
	Attachments string
	Steps       string
	Sourcecode  string
	Runtime     string
	RuntimeVersion string
	Published   string
	ResolutionW string
	ResolutionH string
}

const (
	ProgramID uint = 1 << iota
	ProgramTitle
	ProgramUserID
	ProgramThumbnail
	ProgramDescription
	ProgramStartax
	ProgramAttachments
	ProgramSteps
	ProgramSourcecode
	ProgramRuntime
	ProgramRuntimeVersion
	ProgramPublished
	ProgramResolution
)

func (this *RawProgram) Validate(flag uint) error {

	published := true

	if (flag & ProgramPublished) != 0 {

		if this.Published != "true" {
			published = false
		}

	}

	if (flag & ProgramID) != 0 {

		programID, err := strconv.Atoi(this.ID)
		if err != nil {
			return errors.New("プログラムIDが不正です。")
		}

		if programID < 0 {
			return errors.New("プログラムIDが不正です。")
		}

	}

	if (flag & ProgramTitle) != 0 {

		if len(this.Title) <= 0 || len(this.Title) >= 100 {
			return errors.New("タイトルの文字数が範囲外です。")
		}

	}

	if (flag & ProgramUserID) != 0 {

		// TOO: implement

	}

	if (flag & ProgramThumbnail) != 0 {

		// TODO: implement

	}

	if (flag & ProgramDescription) != 0 {

		if published {
			if len(this.Description) <= 0 || len(this.Description) > 1000 {
				return errors.New("説明文の文字数が範囲外です。")
			}
		} else {
			if len(this.Description) > 1000 {
				return errors.New("説明文の文字数が範囲外です。")
			}
		}

	}

	if (flag & ProgramStartax) != 0 {

		// TODO: implement

	}

	if (flag & ProgramAttachments) != 0 {

		// TODO: implement

	}

	if (flag & ProgramSteps) != 0 {

		steps, err := strconv.Atoi(this.Steps)
		if err != nil {
			return errors.New("ステップ上限数が正常な値ではありません。")
		}

		if 0 <= steps && steps <= 30000 {
		} else {
			return errors.New("ステップ上限数が範囲外です。")
		}
	}

	if (flag & ProgramSourcecode) != 0 {

		// TODO: implement

	}

	if (flag & ProgramResolution) != 0 {

	}

	if (flag & ProgramRuntime) != 0 {
		switch this.Runtime {
		case "HSP3Dish":
		case "HGIMG4":
		default:
			return errors.New("ランタイム名が不正です。")
		}
	}

	if (flag & ProgramRuntimeVersion) != 0 {
		if !config.IsValidRuntimeVersion(this.RuntimeVersion) {
			return errors.New("ランタイムバージョンが不正です。")
		}
	}

	return nil
}

func (this *RawProgram) ToProgram(flag uint) (*Program, error) {

	program := NewProgram()

	program.Published = true

	if (flag & ProgramPublished) != 0 {

		if this.Published != "true" {
			program.Published = false
		}

	}

	if (flag & ProgramID) != 0 {

		programId, err := strconv.Atoi(this.ID)
		if err != nil {
			return program, err
		}

		program.ID = programId

	}

	if (flag & ProgramTitle) != 0 {

		program.Title = this.Title

	}

	if (flag & ProgramUserID) != 0 {

		userId, err := strconv.Atoi(this.UserID)
		if err != nil {
			return program, err
		}
		program.UserID = userId

	}

	if (flag & ProgramDescription) != 0 {

		program.Description = this.Description

	}

	if (flag & ProgramSteps) != 0 {

		steps, err := strconv.Atoi(this.Steps)
		if err != nil {
			return program, err
		}

		program.Steps = steps
	}

	if (flag & ProgramResolution) != 0 {

		resolutionW, err := strconv.Atoi(this.ResolutionW)
		if err != nil {
			return program, err
		}

		resolutionH, err := strconv.Atoi(this.ResolutionH)
		if err != nil {
			return program, err
		}

		program.ResolutionW = resolutionW
		program.ResolutionH = resolutionH
	}

	if (flag & ProgramRuntime) != 0 {
		program.Runtime = this.Runtime
	}

	if (flag & ProgramRuntimeVersion) != 0 {
		program.RuntimeVersion = this.RuntimeVersion
	}

	if (flag & ProgramStartax) != 0 {

		data, err := base64.StdEncoding.DecodeString(this.Startax)
		if err != nil {
			return program, err
		}

		if len(data) == 0 {
			if program.Published {
				return program, errors.New("Startaxファイルの内容が空です。")
			}
		}

		program.Startax.Data = data

	}

	if (flag & ProgramAttachments) != 0 {

		var pairs []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}

		err := json.Unmarshal([]byte(this.Attachments), &pairs)

		if err != nil {
			if program.Published {
				return program, err
			}
			// TODO: elseの場合、抜け出す処理をすべき
		}

		for _, pair := range pairs {

			var data []byte

			if len(pair.Value) == 0 {
				return program, errors.New("空のファイルが送信されました。")
			}

			data, err = base64.StdEncoding.DecodeString(pair.Value)

			if err != nil {
				return program, err
			}

			program.Attachments = append(program.Attachments, Attachment{
				Name: pair.Name,
				Data: data,
			})

		}

	}

	if (flag & ProgramThumbnail) != 0 {

		data, err := base64.StdEncoding.DecodeString(this.Thumbnail)
		if err != nil {
			return program, err
		}

		if len(data) == 0 {
			if program.Published {
				return program, errors.New("サムネイルの内容が空です。")
			}
		}

		program.Thumbnail.Data = data
	}

	if (flag & ProgramSourcecode) != 0 {

		program.Sourcecode = this.Sourcecode

	}

	return program, nil
}

type ProgramColumn int

const (
	ProgramColId ProgramColumn = iota
	ProgramColTitle
	ProgramColUser
	ProgramColDescription
	ProgramColStartax
	ProgramColAttachments
	ProgramColCreatedAt
	ProgramColUpdatedAt
	ProgramColGood
	ProgramColPlay
	ProgramColThumbnail
	ProgramColSteps
	ProgramColSourcecode
	ProgramColResolutionW
	ProgramColResolutionH
	ProgramColRuntime
	ProgramColRuntimeVersion
)

func (this *ProgramColumn) String() string {
	switch *this {
	case ProgramColId:
		return "id"
	case ProgramColTitle:
		return "title"
	case ProgramColUser:
		return "user"
	case ProgramColDescription:
		return "description"
	case ProgramColStartax:
		return "startax"
	case ProgramColAttachments:
		return "attachments"
	case ProgramColCreatedAt:
		return "created_at"
	case ProgramColUpdatedAt:
		return "updated_at"
	case ProgramColGood:
		return "good"
	case ProgramColPlay:
		return "play"
	case ProgramColThumbnail:
		return "thumbnail"
	case ProgramColSteps:
		return "steps"
	case ProgramColSourcecode:
		return "sourcecode"
	case ProgramColResolutionW:
		return "resolution_w"
	case ProgramColResolutionH:
		return "resolution_h"
	case ProgramColRuntime:
		return "runtime"
	case ProgramColRuntimeVersion:
		return "runtime_version"
	}

	return ""
}

func GetProgramRankingForDay(out *[]Program, from int, number int) (int, error) {

	now := time.Now()
	todayBegin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return getProgramRankingSince(todayBegin, out, from, number)
}

func GetProgramRankingForWeek(out *[]Program, from int, number int) (int, error) {

	now := time.Now()
	thisWeekBegin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -7)

	return getProgramRankingSince(thisWeekBegin, out, from, number)
}

func GetProgramRankingForMonth(out *[]Program, from int, number int) (int, error) {

	now := time.Now()
	thisMonthBegin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, -1, 0)

	return getProgramRankingSince(thisMonthBegin, out, from, number)
}

func getProgramRankingSince(since time.Time, out *[]Program, from int, number int) (int, error) {
	var err error

	var rowCount int
	err = DB.Model(Program{}).Scopes(Published).Where("created_at >= ?", since.Format("2006-1-2")).Count(&rowCount).Error
	if err != nil {
		return 0, err
	}

	err = DB.Model(Program{}).Scopes(Published).Where("created_at >= ?", since.Format("2006-1-2")).Order("good desc, play desc").Limit(number).Offset(from).Find(out).Error
	return rowCount, err
}

func GetProgramRankingForAllTime(out *[]Program, from int, number int) (int, error) {

	return GetProgramListBy(ProgramColGood, out, true, from, number)
}

func GetProgramListBy(keyColumn ProgramColumn, out *[]Program, isDesc bool, from int, number int) (int, error) {
	var err error

	// 並び順
	var order string

	if isDesc {
		order = "DESC"
	} else {
		order = "ASC"
	}

	var rowCount int
	err = DB.Table("programs").Count(&rowCount).Error
	if err != nil {
		return 0, err
	}

	// クエリを発行
	err = DB.Model(Program{}).Scopes(Published).Order(keyColumn.String() + " " + order).Limit(number).Offset(from).Find(out).Error

	return rowCount, err
}

func GetProgramListByQuery(out *[]Program, query string, keyColumn ProgramColumn, isDesc bool, number int, offset int) (int, error) {
	var err error

	// 並び順
	var order string

	if isDesc {
		order = "DESC"
	} else {
		order = "ASC"
	}

	queryMod := "%" + query + "%"

	var rowCount int
	err = DB.Model(Program{}).Scopes(Published).Where("title LIKE ?", queryMod).Count(&rowCount).Error
	if err != nil {
		return 0, err
	}

	// クエリを発行
	err = DB.Model(Program{}).Scopes(Published).Where("title LIKE ?", queryMod).Order(keyColumn.String() + " " + order).Limit(number).Offset(offset).Find(out).Error

	return int(rowCount), err

}

func GetProgramListByUser(keyColumn ProgramColumn, out *[]Program, user int, isDesc bool, from int, number int) (int, error) {
	var err error

	// 並び順
	var order string

	if isDesc {
		order = "DESC"
	} else {
		order = "ASC"
	}

	var rowCount int
	err = DB.Model(Program{}).Where("user_id = ?", user).Scopes(Published).Count(&rowCount).Error
	if err != nil {
		return 0, err
	}

	err = DB.Model(Program{}).Where("user_id = ?", user).Scopes(Published).Order(keyColumn.String() + " " + order).Limit(number).Offset(from).Find(out).Error

	return rowCount, err
}

func GetProgramListRelatedTo(out *[]Program, title string, number int) error {
	var err error

	token := ngram.NewTokenize(3, title)

	repl := strings.NewReplacer("_", "", "[", "", "]", "", "%", "", "'", "", "`", "", "\"", "")
	querys := make([]string, 0)
	for _, t := range token.Tokens() {
		querys = append(querys, "%"+repl.Replace(t.String())+"%")
	}

	statement := DB.Model(Program{}).Scopes(Published)
	cond := "`programs`.title <> '" + repl.Replace(title) + "' AND ("
	for i, query := range querys {
		if i != 0 {
			cond += " OR "
		}
		cond += "`programs`.title LIKE '%" + query + "%'"
	}
	cond += ")"

	statement = statement.Where(cond)

	var rowCount int
	err = statement.Count(&rowCount).Error

	if (err != nil) || (rowCount == 0) {
		return errors.New("関連プログラムが見つかりませんでした。")
	}

	err = statement.Find(out).Error
	return err
}

func ExistsProgram(id int) bool {

	var rowCount int
	err := DB.Model(Program{}).Scopes(Published).Where("id = ?", id).Count(&rowCount).Error

	if err != nil {
		return false
	}

	return rowCount > 0
}

func PlayProgram(id int) error {

	err := DB.Model(Program{}).Where("id = ?", id).Update("play", gorm.Expr("play + 1")).Error

	return err
}
