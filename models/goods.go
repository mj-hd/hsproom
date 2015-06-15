package models

import (
	_ "github.com/go-sql-driver/mysql"
)

func initGoods() {
	DB.AutoMigrate(&Good{})
}

type Good struct {
	ID        int `gorm:"primary_key"`
	UserID    int `sql:"index"`
	ProgramID int `sql:"index"`
}

func (this *Good) Load(id int) error {

	err := DB.Model(Good{}).First(this, id).Error

	return err
}

func (this *Good) LoadByUserAndProgram(userId int, programId int) error {

	err := DB.Model(Good{}).Where("user_id = ? AND program_id", userId, programId).First(this).Error

	return err
}

func (this *Good) Create() (int, error) {

	err := DB.Create(this).Error

	return this.ID, err
}

func (this *Good) Remove() error {

	err := DB.Delete(this).Error

	return err
}

func GetGoodListByUser(out *[]Good, userId int, from int, number int) (int, error) {
	var err error

	if cap(*out) < number {
		*out = make([]Good, number)
	}

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

	if cap(*out) < number {
		*out = make([]Good, number)
	}

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
