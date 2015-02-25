package models

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Program struct {
	*ProgramInfo
	Startax     []byte
	Attachments *Attachments
	Thumbnail   []byte
}

type ProgramInfo struct {
	Id          int
	Created     time.Time
	Modified    mysql.NullTime
	Title       string
	User        string
	UserId      int
	Good        int
	Play        int
	Description string
	Size        int
}

type Attachments struct {
	Files []File
}

type File struct {
	Name string
	Data []byte
}

func NewProgram() *Program {
	return &Program{
		ProgramInfo: &ProgramInfo{},
		Startax:     make([]byte, 0),
		Attachments: &Attachments{
			Files: make([]File, 0),
		},
	}
}

func (this *Program) Load(id int) error {

	var rawAttachments []byte

	row := DB.QueryRow("SELECT id, created, modified, title, user, good, play, thumbnail, description, startax, size, attachments FROM programs WHERE id = ?", id)
	err := row.Scan(&this.Id, &this.Created, &this.Modified, &this.Title, &this.User, &this.Good, &this.Play, &this.Thumbnail, &this.Description, &this.Startax, &this.Size, &rawAttachments)

	if err != nil {
		return err
	}

	if rawAttachments == nil {
		return nil
	}

	buffer := bytes.NewBuffer(rawAttachments)
	decoder := gob.NewDecoder(buffer)

	err = decoder.Decode(&this.Attachments)

	if err != nil {
		return err
	}

	this.UserId, err = GetUserIdFromName(this.User)

	return err
}

func (this *Program) Update() error {

	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)

	err := encoder.Encode(this.Attachments)
	if err != nil {
		return err
	}

	_, err = DB.Exec("UPDATE programs SET modified = ?, title = ?, thumbnail = ?, description = ?, startax = ?, size = ?, attachments = ? WHERE id = ?",
		time.Now(), this.Title, this.Thumbnail, this.Description, this.Startax, this.Size, buffer.Bytes(), this.Id)

	if err != nil {
		return err
	}

	return nil
}

func (this *Program) Create() (int, error) {

	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)

	err := encoder.Encode(this.Attachments)
	if err != nil {
		return 0, err
	}

	result, err := DB.Exec("INSERT INTO programs ( created, title, user, thumbnail, description, startax, size, attachments ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ? )", time.Now(), this.Title, this.User, this.Thumbnail, this.Description, this.Startax, this.Size, buffer.Bytes())
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func (this *Program) Remove() error {

	_, err := DB.Exec("DELETE FROM programs WHERE id = ?", this.Id)

	return err
}

func (this *ProgramInfo) Load(id int) error {

	row := DB.QueryRow("SELECT id, created, modified, title, user, good, play, description, size FROM programs WHERE id = ?", id)
	err := row.Scan(&this.Id, &this.Created, &this.Modified, &this.Title, &this.User, &this.Good, &this.Play, &this.Description, &this.Size)

	if err != nil {
		return err
	}

	this.UserId, err = GetUserIdFromName(this.User)
	if err != nil {
		return err
	}

	return nil
}

func (this *ProgramInfo) Update() error {

	_, err := DB.Exec("UPDATE programs SET modified = ?, title = ?, description = ? WHERE id = ?",
		time.Now(), this.Title, this.Description, this.Id)

	if err != nil {
		return err
	}

	return nil
}

func (this *ProgramInfo) Remove() error {

	_, err := DB.Exec("DELETE FROM programs WHERE id = ?", this.Id)

	return err
}

type RawProgram struct {
	Id          string
	Created     string
	Modified    string
	Title       string
	User        string
	UserId      string
	Thumbnail   string
	Description string
	Startax     string
	Size        string
	Attachments string
}

const (
	ProgramId uint = 1 << iota
	ProgramCreated
	ProgramModified
	ProgramTitle
	ProgramUser
	ProgramUserId
	ProgramThumbnail
	ProgramDescription
	ProgramStartax
	ProgramSize
	ProgramAttachments
)

func (this *RawProgram) Validate(flag uint) error {

	if (flag & ProgramId) != 0 {

		programId, err := strconv.Atoi(this.Id)
		if err != nil {
			return errors.New("プログラムIDが不正です。")
		}

		if programId < 0 {
			return errors.New("プログラムIDが不正です。")
		}

	}

	if (flag & ProgramCreated) != 0 {

		// TODO: implement

	}

	if (flag & ProgramModified) != 0 {

		// TODO: implement

	}

	if (flag & ProgramTitle) != 0 {

		if len(this.Title) <= 0 || len(this.Title) >= 100 {
			return errors.New("タイトルの文字数が範囲外です。")
		}

	}

	if (flag & ProgramUser) != 0 {

		if len(this.User) <= 0 || len(this.User) >= 50 {
			return errors.New("ユーザ名の文字数が範囲外です。")
		}

	}

	if (flag & ProgramUserId) != 0 {

		// TODO: implement

	}

	if (flag & ProgramThumbnail) != 0 {

		// TODO: implement

	}

	if (flag & ProgramDescription) != 0 {

		if len(this.Description) <= 0 {
			return errors.New("説明文の文字数が範囲外です。")
		}

	}

	if (flag & ProgramStartax) != 0 {

		// TODO: implement

	}

	if (flag & ProgramSize) != 0 {

		// TODO: implement

	}

	if (flag & ProgramAttachments) != 0 {

		// TODO: implement

	}

	return nil
}

func (this *RawProgram) ToProgram(flag uint) (*Program, error) {

	program := NewProgram()
	oldProgram := NewProgram()

	programInfo, err := this.ToProgramInfo(flag)
	if err != nil {
		return program, err
	}

	program.ProgramInfo = &programInfo

	if (flag & ProgramStartax) != 0 {

		data, err := base64.StdEncoding.DecodeString(this.Startax)
		if err != nil {
			return program, err
		}

		if len(data) == 0 {
			if program.Id == 0 {
				return program, errors.New("Startaxファイルの内容が空です。")
			}

			if oldProgram.Id == 0 {
				err = oldProgram.Load(program.Id)

				if err != nil {
					return program, errors.New("内部エラーが発生しました。")
				}
			}

			data = oldProgram.Startax
		}

		program.Startax = data

	}

	if (flag & ProgramAttachments) != 0 {

		var pairs []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}

		err := json.Unmarshal([]byte(this.Attachments), &pairs)

		if err != nil {
			return program, err
		}

		for _, pair := range pairs {

			var data []byte

			if pair.Value == "PASS" {

				if oldProgram.Id == 0 {
					err = oldProgram.Load(program.Id)

					if err != nil {
						return program, errors.New("内部エラーが発生しました。")
					}
				}

				var file File

				for _, file = range oldProgram.Attachments.Files {
					if pair.Name == file.Name {
						break
					}
				}

				program.Attachments.Files = append(program.Attachments.Files, file)

			} else if pair.Value == "DELETE" {
			} else {

				data, err = base64.StdEncoding.DecodeString(pair.Value)

				if err != nil {
					return program, err
				}
			}

			program.Attachments.Files = append(program.Attachments.Files, File{
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
			if program.Id == 0 {
				return program, errors.New("サムネイルの内容が空です。")
			}

			if oldProgram.Id == 0 {
				err = oldProgram.Load(program.Id)

				if err != nil {
					return program, errors.New("内部エラーが発生しました。")
				}
			}

			data = oldProgram.Thumbnail
		}

		program.Thumbnail = data
	}

	return program, nil
}

func (this *RawProgram) ToProgramInfo(flag uint) (ProgramInfo, error) {

	var program ProgramInfo

	if (flag & ProgramId) != 0 {

		programId, err := strconv.Atoi(this.Id)
		if err != nil {
			return program, err
		}

		program.Id = programId

	}

	if (flag & ProgramCreated) != 0 {

		// TODO: implement

	}

	if (flag & ProgramModified) != 0 {

		// TODO: implement

	}

	if (flag & ProgramTitle) != 0 {

		program.Title = this.Title

	}

	if (flag & ProgramUser) != 0 {

		program.User = this.User

	}

	if (flag & ProgramUserId) != 0 {

		id, err := strconv.Atoi(this.UserId)

		if err != nil {
			return program, err
		}

		program.UserId = id

	}

	if (flag & ProgramDescription) != 0 {

		program.Description = this.Description

	}

	if (flag & ProgramSize) != 0 {

		size, err := strconv.Atoi(this.Size)
		if err != nil {
			return program, err
		}

		program.Size = size

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
	ProgramColCreated
	ProgramColModified
	ProgramColGood
	ProgramColPlay
	ProgramColThumbnail
	ProgramColSize
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
	case ProgramColCreated:
		return "created"
	case ProgramColModified:
		return "modified"
	case ProgramColGood:
		return "good"
	case ProgramColPlay:
		return "play"
	case ProgramColThumbnail:
		return "thumbnail"
	case ProgramColSize:
		return "size"
	}

	return ""
}

func GetProgramRankingForDay(out *[]ProgramInfo, from int, number int) (int, error) {

	now := time.Now()
	todayBegin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return getProgramRankingSince(todayBegin, out, from, number)
}

func GetProgramRankingForWeek(out *[]ProgramInfo, from int, number int) (int, error) {

	now := time.Now()
	thisWeekBegin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -7)

	return getProgramRankingSince(thisWeekBegin, out, from, number)
}

func GetProgramRankingForMonth(out *[]ProgramInfo, from int, number int) (int, error) {

	now := time.Now()
	thisMonthBegin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, -1, 0)

	return getProgramRankingSince(thisMonthBegin, out, from, number)
}

func getProgramRankingSince(since time.Time, out *[]ProgramInfo, from int, number int) (int, error) {

	// キャパシティチェック
	if cap(*out) < number {
		*out = make([]ProgramInfo, number)
	}

	var rowCount int
	err := DB.QueryRow("SELECT count(id) FROM programs WHERE created >= ?", since.Format("2006-1-2")).Scan(&rowCount)

	if err != nil {
		return 0, err
	}

	rows, err := DB.Query("SELECT id FROM programs WHERE created >= ? ORDER BY good DESC LIMIT ?, ?", since.Format("2006-1-2"), from, number)

	if err != nil {
		return rowCount, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var id int
		err := rows.Scan(&id)

		if err != nil {
			return rowCount, err
		}

		err = (*out)[i].Load(id)

		if err != nil {
			return rowCount, err
		}

		i++
	}

	return rowCount, nil

}

func GetProgramRankingForAllTime(out *[]ProgramInfo, from int, number int) (int, error) {

	return GetProgramListBy(ProgramColGood, out, true, from, number)
}

func GetProgramListBy(keyColumn ProgramColumn, out *[]ProgramInfo, isDesc bool, from int, number int) (int, error) {

	// キャパシティチェック
	if cap(*out) < number {
		*out = make([]ProgramInfo, number)
	}

	// 並び順
	var order string

	if isDesc {
		order = "DESC"
	} else {
		order = "ASC"
	}

	// クエリを発行
	rows, err := DB.Query("SELECT id FROM programs ORDER BY "+keyColumn.String()+" "+order+" LIMIT ?, ?", from, number)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	// outへ格納
	i := 0
	for rows.Next() {
		var id int
		err := rows.Scan(&id)

		if err != nil {
			return i, err
		}

		err = (*out)[i].Load(id)

		if err != nil {
			return i, err
		}

		i++
	}

	return i, nil
}

func GetProgramListByQuery(out *[]ProgramInfo, query string, keyColumn ProgramColumn, isDesc bool, number int, offset int) (int, error) {

	if cap(*out) < number {
		*out = make([]ProgramInfo, number)
	}

	// 並び順
	var order string

	if isDesc {
		order = "DESC"
	} else {
		order = "ASC"
	}

	queryMod := "%" + query + "%"

	rowCount := 0

	err := DB.QueryRow("SELECT count(*) FROM programs WHERE title LIKE ?", queryMod).Scan(&rowCount)

	if err != nil {
		return rowCount, err
	}

	// クエリを発行
	rows, err := DB.Query("SELECT id FROM programs WHERE title LIKE ? ORDER BY "+keyColumn.String()+" "+order+" LIMIT ?, ?", queryMod, offset, number)
	if err != nil {
		return rowCount, err
	}
	defer rows.Close()

	// outへ格納
	i := 0
	for rows.Next() {
		var id int
		err := rows.Scan(&id)

		if err != nil {
			return rowCount, err
		}

		err = (*out)[i].Load(id)

		if err != nil {
			return rowCount, err
		}

		i++
	}

	return rowCount, nil

}

func GetProgramListByUser(keyColumn ProgramColumn, out *[]ProgramInfo, name string, isDesc bool, from int, number int) (int, error) {

	if cap(*out) < number {
		*out = make([]ProgramInfo, number)
	}

	// 並び順
	var order string

	if isDesc {
		order = "DESC"
	} else {
		order = "ASC"
	}

	var rowCount int
	err := DB.QueryRow("SELECT count(id) FROM programs WHERE user = ?", name).Scan(&rowCount)

	if err != nil {
		return 0, err
	}

	// クエリを発行
	rows, err := DB.Query("SELECT id FROM programs WHERE user = ? ORDER BY "+keyColumn.String()+" "+order+" LIMIT ?, ?", name, from, number)
	if err != nil {
		return rowCount, err
	}
	defer rows.Close()

	// outへ格納
	i := 0
	for rows.Next() {
		var id int
		err := rows.Scan(&id)

		if err != nil {
			return rowCount, err
		}

		err = (*out)[i].Load(id)

		if err != nil {
			return rowCount, err
		}

		i++
	}

	return rowCount, nil
}

func ExistsProgram(id int) bool {

	var rowCount int
	err := DB.QueryRow("SELECT count(id) FROM programs WHERE id = ?", id).Scan(&rowCount)

	if err != nil {
		return false
	}

	if rowCount < 1 {
		return false
	}

	return true
}

func PlayProgram(id int) error {

	_, err := DB.Exec("UPDATE programs SET play = play + 1 WHERE id = ?", id)

	return err
}
