package models

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func initComments() {
	DB.AutoMigrate(&Comment{})
}

type Comment struct {
	ID        int
	CreatedAt time.Time  ``
	DeletedAt *time.Time ``
	Message   string     `sql:"size:200;not null;"`
	Program   Program
	ProgramID int `sql:"index"`
	User      User
	UserID    int    `sql:"index"`
	UserName  string `sql:"size:100"`
	ReplyTo   int
	Replies   []Comment
}

func (this *Comment) Load(id int) error {
	return DB.First(this, id).Error
}

func (this *Comment) LoadReplies() error {
	return DB.Model(this).Where("reply_to = ?", this.ID).Find(&this.Replies).Error
}

func (this *Comment) LoadProgram() error {
	return DB.First(&this.Program, this.ProgramID).Error
}

func (this *Comment) LoadUser() error {
	return DB.First(&this.User, this.UserID).Error
}

func (this *Comment) Update() error {
	return DB.Save(this).Error
}

func (this *Comment) Create() error {
	return DB.Create(this).Error
}

func (this *Comment) Remove() error {
	return DB.Delete(this).Error
}

func GetComments(programId int, number int, offset int, since int) (result []Comment, err error) {
	result = make([]Comment, number)
	err = DB.Model(Comment{}).Where("id > ? and program_id = ? and reply_to = -1", since, programId).Order("created_at desc").Limit(number).Offset(offset).Find(&result).Error
	return result, err
}

func GetCommentsAndReplies(programId int, number int, offset int, since int) (result []Comment, err error) {
	result = make([]Comment, number)
	err = DB.Model(Comment{}).Where("id > ? and program_id = ?", since, programId).Order("created_at desc").Limit(number).Offset(offset).Find(&result).Error
	return result, err
}

func GetCommentsCount(programId int) (tot int, err error) {
	err = DB.Model(Comment{}).Where("program_id = ? and reply_to = -1", programId).Count(&tot).Error
	return tot, err
}

func GetCommentsAndRepliesMaxID(programId int) (max int, err error) {
	var maxes []int
	err = DB.Model(Comment{}).Where("program_id = ?", programId).Order("created_at desc").Limit(1).Pluck("id", &maxes).Error
	return maxes[0], err
}
