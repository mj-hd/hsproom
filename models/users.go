package models

type User struct {
	Id       int
	Name     string
	ScreenName string
	Token    string
	Secret   string
	Profile  string
	IconURL  string
	Website  string
	Location string
}

func (this *User) Load(id int) error {

	row := DB.QueryRow("SELECT id, name, screenname, profile, icon_url, token, secret, website, location FROM users WHERE id = ?", id)
	err := row.Scan(&this.Id, &this.Name, &this.ScreenName, &this.Profile, &this.IconURL, &this.Token, &this.Secret, &this.Website, &this.Location)

	if err != nil {
		return err
	}

	return nil
}

func (this *User) LoadFromScreenName(screenname string) error {
	row := DB.QueryRow("SELECT id, name, screenname, profile, icon_url, token, secret, website, location FROM users WHERE screenname = ?", screenname)
	err := row.Scan(&this.Id, &this.Name, &this.ScreenName, &this.Profile, &this.IconURL, &this.Token, &this.Secret, &this.Website, &this.Location)

	if err != nil {
		return err
	}

	return nil
}

func (this *User) Update() error {

	_, err := DB.Exec("UPDATE users SET name = ?, screenname = ?, profile = ?, icon_url = ?, token = ?, secret = ?, website = ?, location = ? WHERE id = ?",
		this.Name, this.ScreenName, this.Profile, this.IconURL, this.Token, this.Secret, this.Website, this.Location, this.Id)

	if err != nil {
		return err
	}

	return nil
}

func (this *User) Create() (int, error) {

	result, err := DB.Exec("INSERT INTO users ( name, screenname, profile, icon_url, token, secret, website, location) VALUES ( ?, ?, ?, ?, ?, ?, ?, ? )", this.Name, this.ScreenName, this.Profile, this.IconURL, this.Token, this.Secret, this.Website, this.Location)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func (this *User) Remove() error {

	_, err := DB.Exec("DELETE FROM users WHERE id = ?", this.Id)

	return err
}

func ExistsUserScreenName(screenname string) bool {

	var rowCount int
	err := DB.QueryRow("SELECT count(id) FROM users WHERE screenname = ?", screenname).Scan(&rowCount)

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
	err := DB.QueryRow("SELECT count(id) FROM users WHERE id = ?", id).Scan(&rowCount)

	if err != nil {
		return false
	}

	if rowCount < 1 {
		return false
	}

	return true
}

func GetUserName(id int) (string, error) {

	var name string

	row := DB.QueryRow("SELECT name FROM users WHERE id = ?", id)
	err := row.Scan(&name)

	if err != nil {
		return "", err
	}

	return name, nil
}

func GetUserScreenName(id int) (string, error) {

	var name string

	row := DB.QueryRow("SELECT screenname FROM users WHERE id = ?", id)
	err := row.Scan(&name)

	if err != nil {
		return "", err
	}

	return name, nil
}

func GetUserIdFromScreenName(screenname string) (int, error) {

	var id int

	row := DB.QueryRow("SELECT id FROM users WHERE screenname = ?", screenname)
	err := row.Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
