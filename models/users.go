package models

import (
	"time"

	"../config"
)

func initUsers() {
	DB.AutoMigrate(&User{})
}

type User struct {
	ID         int
	CreatedAt  time.Time  ``
	UpdatedAt  time.Time  ``
	DeletedAt  *time.Time ``
	Name       string     `sql:"size:50;not null;"`
	ScreenName string     `sql:"size:50;not null;"`
	Profile    string     `sql:"size:300;default:''"`
	IconURL    string     `sql:"size:140;default:''"`
	Website    string     `sql:"size:300;default:''"`
	Location   string     `sql:"size:50;default:''"`

	Programs []Program
	Goods    []Good
	Notifications []Notification
}

func (this *User) AfterFind() (err error) {
	this.CreatedAt = this.CreatedAt.In(config.JST())
	this.UpdatedAt = this.UpdatedAt.In(config.JST())

	if this.DeletedAt != nil {
		*this.DeletedAt = this.DeletedAt.In(config.JST())
	}

	return nil
}

func (this *User) Load(id int) error {

	err := DB.First(this, id).Error

	return err
}

func (this *User) LoadFromScreenName(screenname string) error {

	err := DB.Where("screen_name = ?", screenname).First(this).Error

	return err
}

func (this *User) LoadPrograms() error {
	return DB.Model(this).Related(&this.Programs).Error
}

func (this *User) LoadGoods() error {
	return DB.Model(this).Related(&this.Goods).Error
}

func (this *User) LoadNotifications() error {
	return DB.Model(this).Related(&this.Notifications).Error
}


func (this *User) Update() error {

	err := DB.Save(this).Error

	return err
}

func (this *User) Create() (int, error) {

	err := DB.Create(this).Error

	return this.ID, err
}

func (this *User) Remove() error {

	err := DB.Delete(this).Error

	return err
}

func ExistsUserScreenName(screenname string) bool {

	var rowCount int
	err := DB.Model(User{}).Where("screen_name = ?", screenname).Count(&rowCount).Error

	if err != nil {
		return false
	}

	if rowCount < 1 {
		return false
	}

	return true
}

func ExistsUser(id int) bool {

	var rowCount int
	err := DB.Model(User{}).Where("id = ?", id).Count(&rowCount).Error

	if err != nil {
		return false
	}

	if rowCount < 1 {
		return false
	}

	return true
}

func GetUserName(id int) (string, error) {

	result := struct {
		Name string
	}{
		Name: "",
	}

	err := DB.Model(User{}).Where("id = ?", id).Select("name").Scan(&result).Error

	return result.Name, err
}

func GetUserScreenName(id int) (string, error) {

	result := struct {
		Name string
	}{
		Name: "",
	}

	err := DB.Model(User{}).Where("id = ?", id).Select("screen_name").Scan(&result).Error

	return result.Name, err
}

func GetUserIdFromScreenName(screenname string) (int, error) {

	var id int

	err := DB.Model(User{}).Where("screen_name = ?", screenname).Select("id").Scan(&id).Error

	return id, err
}
