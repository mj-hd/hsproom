package models

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"../config"
)

func initGoods() {
	DB.AutoMigrate(&Good{})
}

type Good struct {
	ID        int `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	UserID    int `sql:"index"`
	ProgramID int `sql:"index"`
}

func (this *Good) AfterFind() (err error) {
	this.CreatedAt = this.CreatedAt.In(config.JST())
	this.UpdatedAt = this.UpdatedAt.In(config.JST())

	if this.DeletedAt != nil {
		*this.DeletedAt = this.DeletedAt.In(config.JST())
	}

	return nil
}

func (this *Good) Load(id int) error {

	err := DB.Model(Good{}).First(this, id).Error

	return err
}

func (this *Good) LoadByUserAndProgram(userId int, programId int) error {

	err := DB.Model(Good{}).Where("user_id = ? AND program_id = ?", userId, programId).First(this).Error

	return err
}

func (this *Good) Create() (int, error) {

	tx := DB.Begin()

	err := tx.Create(this).Error
	if err != nil {
		tx.Rollback()
		return this.ID, err
	}

	err = tx.Model(Program{}).Where("id = ?", this.ProgramID).UpdateColumn("good", gorm.Expr("good + ?", 1)).Error
	if err != nil {
		tx.Rollback()
		return this.ID, err
	}

	tx.Commit()

	return this.ID, nil
}

func (this *Good) Remove() error {

	tx := DB.Begin()

	err := tx.Delete(this).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = DB.Model(Program{}).Where("id = ?", this.ProgramID).UpdateColumn("good", gorm.Expr("good - ?", 1)).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return err
}

func GetGoodListByUser(out *[]Good, userId int, from int, number int) (int, error) {
	var err error

	var rowCount int
	err = DB.Model(Good{}).Where("user_id = ?", userId).Count(&rowCount).Error

	if err != nil {
		return 0, err
	}

	err = DB.Where("user_id = ?", userId).Limit(number).Offset(from).Find(out).Error

	return rowCount, err
}

func GetGoodListByProgram(out *[]Good, programId int, from int, number int) (int, error) {
	var err error

	var rowCount int
	err = DB.Model(Good{}).Where("program_id = ?", programId).Count(&rowCount).Error

	if err != nil {
		return 0, err
	}

	err = DB.Where("program_id = ?", programId).Limit(number).Offset(from).Find(out).Error

	return rowCount, err
}

func CanGoodProgram(userId int, programId int) bool {

	var rowCount int

	err := DB.Model(Good{}).Where("user_id = ? AND program_id = ?", userId, programId).Count(&rowCount).Error

	if err != nil {
		return false
	}

	return rowCount < 1
}

func GetGoodCountByProgram(programId int) int {

	var rowCount int
	err := DB.Model(Good{}).Where("program_id = ?", programId).Count(&rowCount).Error

	if err != nil {
		return 0
	}

	return rowCount
}
