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
	ProgramID int
	User      User
	UserID    int
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
	return DB.Model(this).Related(&this.Program).Error
}

func (this *Comment) LoadUser() error {
	return DB.Model(this).Related(&this.User).Error
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

func GetComments(programId int, number int, offset int) (result []Comment, err error) {
	err = DB.Model(Comment{}).Scopes(Published).Order("created_at desc").Limit(number).Offset(offset).Find(result).Error
	return result, err
}
