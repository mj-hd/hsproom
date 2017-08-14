package models

import (
	"time"

	"../config"

	_ "github.com/go-sql-driver/mysql"
)

func initNotifications() {
	DB.AutoMigrate(&Notification{})
}

type Notification struct {
	ID int
	CreatedAt time.Time ``
	DeletedAt *time.Time ``
	Message string `sql:"size: 1000; not null;"`
	User User
	UserID int `sql:"index"`
	URL string `sql:"size: 500;"`
}

type NotificationColumn int

const (
	NotificationColId NotificationColumn = iota
	NotificationColCreatedAt
	NotificationColDeletedAt
	NotificationColMessage
	NotificationColUser
	NotificationColURL
)

func (this *NotificationColumn) String() string {
	switch *this {
	case NotificationColId:
		return "id"
	case NotificationColCreatedAt:
		return "created_at"
	case NotificationColDeletedAt:
		return "deleted_at"
	case NotificationColMessage:
		return "message"
	case NotificationColUser:
		return "user"
	case NotificationColURL:
		return "url"
	}

	return ""
}

func (this *Notification) Load(id int) error {
	return DB.First(this, id).Error
}

func (this *Notification) LoadUser() error {
	return DB.First(&this.User, this.UserID).Error
}

func (this *Notification) Update() error {
	return DB.Save(this).Error
}

func (this *Notification) Create() error {
	return DB.Create(this).Error
}

func (this *Notification) Remove() error {
	return DB.Delete(this).Error
}

func (this *Notification) AfterFind() error {
	this.CreatedAt = this.CreatedAt.In(config.JST())

	if this.DeletedAt != nil {
		*this.DeletedAt = this.DeletedAt.In(config.JST())
	}

	return nil
}

func GetNotificationListByUser(keyColumn NotificationColumn, out *[]Notification, user int, isDesc bool, from int, number int) (int, error) {
	var err error

	// 並び順
	var order string

	if isDesc {
		order = "DESC"
	} else {
		order = "ASC"
	}

	var rowCount int
	err = DB.Model(Notification{}).Where("user_id = ?", user).Count(&rowCount).Error
	if err != nil {
		return 0, err
	}

	if number == 0 {
		number = rowCount
	}

	err = DB.Model(Notification{}).Where("user_id = ?", user).Order(keyColumn.String() + " " + order).Limit(number).Offset(from).Find(out).Error

	return rowCount, err
}
