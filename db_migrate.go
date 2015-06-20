package main

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"

	"./config"
)

type Program struct {
	*ProgramInfo
	Startax     []byte
	Attachments *Attachments
	Thumbnail   []byte
	Sourcecode  string
}

type ProgramInfo struct {
	Id          int
	Created     time.Time
	Modified    mysql.NullTime
	Title       string
	User        int
	UserName    string
	Good        int
	Play        int
	Description string
	Steps       int
	Runtime     string
}

type Attachments struct {
	Files []File
}

type File struct {
	Name string
	Data []byte
}

type User struct {
	Id         int
	Name       string
	ScreenName string
	Profile    string
	IconURL    string
	Website    string
	Location   string
}

type Good struct {
	Id      int
	User    int
	Program int
}

func main() {
	var err error

	fmt.Println("DBの移行を開始します。")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	migratePrograms(tx)
	migrateUsers(tx)
	migrateGoods(tx)

	cleanPrograms(tx)
	cleanUsers(tx)
	cleanGoods(tx)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func migratePrograms(tx *sql.Tx) {
	fmt.Println("テーブル programs:")

	programs := make([]Program, 0)

	rows, err := tx.Query("SELECT id, created, modified, title, user, good, play, thumbnail, description, startax, attachments, steps, sourcecode, runtime FROM programs")

	if err != nil {
		tx.Rollback()
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		program := Program{
			ProgramInfo: &ProgramInfo{},
			Startax:     make([]byte, 0),
			Attachments: &Attachments{
				Files: make([]File, 0),
			},
		}
		var buffer []byte
		err = rows.Scan(&program.Id, &program.Created, &program.Modified, &program.Title, &program.User, &program.Good, &program.Play, &program.Thumbnail, &program.Description, &program.Startax, &buffer, &program.Steps, &program.Sourcecode, &program.Runtime)

		if err != nil {
			tx.Rollback()
			panic(err)
		}

		fmt.Println("program %d を読み込み中...", program.Id)

		byteBuffer := bytes.NewBuffer(buffer)
		decoder := gob.NewDecoder(byteBuffer)

		err = decoder.Decode(&program.Attachments)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		programs = append(programs, program)
	}

	for _, program := range programs {
		fmt.Println("program %d を更新中...", program.Id)

		for _, att := range program.Attachments.Files {
			_, err = tx.Exec("INSERT INTO attachments (created_at, updated_at, deleted_at, program_id, name, data) VALUES (?, ?, ?, ?, ?, ?)", program.Created, program.Modified.Time, nil, program.Id, att.Name, att.Data)

			if err != nil {
				tx.Rollback()
				panic(err)
			}
		}

		_, err = tx.Exec("INSERT INTO startaxes (created_at, updated_at, deleted_at, program_id, data) VALUES (?, ?, ?, ?, ?)", program.Created, program.Modified.Time, nil, program.Id, program.Startax)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
		_, err = tx.Exec("INSERT INTO thumbnails (created_at, updated_at, deleted_at, program_id, data) VALUES (?, ?, ?, ?, ?)", program.Created, program.Modified.Time, nil, program.Id, program.Thumbnail)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		_, err = tx.Exec("UPDATE programs SET created_at = ?, updated_at = ?, deleted_at = ?, title = ?, user_id = ?, play = ?, description = ?, steps = ?, runtime = ?, published = ?, sourcecode = ? WHERE id = ?", program.Created, program.Modified.Time, nil, program.Title, program.User, program.Play, program.Description, program.Steps, program.Runtime, true, program.Sourcecode, program.Id)

		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	fmt.Println("テーブル programs の更新が終わりました。")

}

func migrateUsers(tx *sql.Tx) {
	fmt.Println("テーブル users:")

	users := make([]User, 0)

	rows, err := tx.Query("SELECT id, name, screenname, profile, icon_url, website, location FROM users")
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Name, &user.ScreenName, &user.Profile, &user.IconURL, &user.Website, &user.Location)

		if err != nil {
			tx.Rollback()
			panic(err)
		}

		fmt.Println("user %d を読み込み中...", user.Id)

		users = append(users, user)
	}

	for _, user := range users {
		fmt.Println("user %d を更新中...", user.Id)

		_, err = tx.Exec("UPDATE users SET created_at = ?, updated_at = ?, deleted_at = ?, name = ?, screen_name = ?, profile = ?, icon_url = ?, website = ?, location = ? WHERE id = ?", time.Now(), time.Now(), nil, user.Name, user.ScreenName, user.Profile, user.IconURL, user.Website, user.Location, user.Id)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	fmt.Println("テーブル users の更新が終わりました。")
}

func migrateGoods(tx *sql.Tx) {
	fmt.Println("テーブル goods:")

	goods := make([]Good, 0)

	rows, err := tx.Query("SELECT id, user, program FROM goods")
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var good Good

		err = rows.Scan(&good.Id, &good.User, &good.Program)

		if err != nil {
			tx.Rollback()
			panic(err)
		}

		fmt.Println("good %d を読み込んでいます...", good.Id)

		goods = append(goods, good)
	}

	for _, good := range goods {

		fmt.Println("good %d を更新しています...", good.Id)

		_, err = tx.Exec("UPDATE goods SET created_at = ?, updated_at = ?, deleted_at = ?, user_id = ?, program_id = ? WHERE id = ?", time.Now(), time.Now(), nil, good.User, good.Program, good.Id)

		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	fmt.Println("テーブル goods の更新が終わりました。")
}

func cleanPrograms(tx *sql.Tx) {
	fmt.Println("テーブル programs:")

	_, err := tx.Exec("ALTER TABLE programs DROP created, DROP modified, DROP user, DROP thumbnail, DROP startax, DROP attachments")
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	fmt.Println("テーブル programs を綺麗にしました。")
}

func cleanUsers(tx *sql.Tx) {
	fmt.Println("テーブル users:")

	_, err := tx.Exec("ALTER TABLE users DROP screenname")
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	fmt.Println("テーブル users を綺麗にしました。")
}

func cleanGoods(tx *sql.Tx) {
	fmt.Println("テーブル goods:")

	_, err := tx.Exec("DROP TRIGGER good_count_increment")
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = tx.Exec("DROP TRIGGER good_count_decrement")
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = tx.Exec("ALTER TABLE goods DROP FOREIGN KEY goods_ibfk_5, DROP FOREIGN KEY goods_ibfk_6")
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = tx.Exec("ALTER TABLE goods DROP user, DROP program")
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	fmt.Println("テーブル goods を綺麗にしました。")
}
