package models

type Good struct {
	Id      int
	User    int
	Program int
}

func (this *Good) Load(id int) error {

	row := DB.QueryRow("SELECT id, user, program FROM goods WHERE id = ?", id)
	err := row.Scan(&this.Id, &this.User, &this.Program)

	if err != nil {
		return err
	}

	return nil
}

func (this *Good) Create() (int, error) {

	result, err := DB.Exec("INSERT INTO goods ( user, program ) VALUES ( ?, ? )", this.User, this.Program)

	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	// TODO: int64をintにダウンキャスト。大丈夫だろうか
	//       users.goにも同様の記述あり
	return int(id), nil
}

func (this *Good) Remove() error {

	_, err := DB.Exec("DELETE FROM goods WHERE id = ?", this.Id)

	return err
}

func GetGoodListByUser(out *[]Good, userId int, from int, number int) (int, error) {

	if cap(*out) < number {
		*out = make([]Good, number)
	}

	var rowCount int
	err := DB.QueryRow("SELECT count(id) FROM goods WHERE user = ?", userId).Scan(&rowCount)

	if err != nil {
		return 0, err
	}

	rows, err := DB.Query("SELECT id FROM goods WHERE user = ? LIMIT ?,?", userId, from, number)

	if err != nil {
		return rowCount, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {

		var id int
		err = rows.Scan(&id)

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

func GetGoodListByProgram(out *[]Good, programId int, from int, number int) (int, error) {

	if cap(*out) < number {
		*out = make([]Good, number)
	}

	var rowCount int
	err := DB.QueryRow("SELECT count(id) FROM goods WHERE program = ?", programId).Scan(&rowCount)

	if err != nil {
		return 0, err
	}

	rows, err := DB.Query("SELECT id FROM goods WHERE program = ? LIMIT ?,?", programId, from, number)

	if err != nil {
		return rowCount, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {

		var id int
		err = rows.Scan(&id)

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

func CanGoodProgram(userId int, programId int) bool {

	var rowCount int

	err := DB.QueryRow("SELECT count(id) FROM goods WHERE user = ? AND program = ?", userId, programId).Scan(&rowCount)

	if err != nil {
		return false
	}

	if rowCount > 0 {
		return false
	}

	return true
}

func GetGoodCountByProgram(programId int) int {

	var rowCount int

	err := DB.QueryRow("SELECT count(id) FROM goods WHERE program = ?", programId).Scan(&rowCount)

	if err != nil {
		return 0
	}

	return rowCount
}
