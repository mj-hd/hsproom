package models

type User struct {
	Id       int
	Name     string
	Token    string
	Secret   string
	Profile  string
	IconURL  string
	Website  string
	Location string
}

func (this *User) Load(id int) error {

	row := DB.QueryRow("SELECT id, name, profile, icon_url, token, secret, website, location FROM users WHERE id = ?", id)
	err := row.Scan(&this.Id, &this.Name, &this.Profile, &this.IconURL, &this.Token, &this.Secret, &this.Website, &this.Location)

	if err != nil {
		return err
	}

	return nil
}

func (this *User) LoadFromName(name string) error {
	row := DB.QueryRow("SELECT id, name, profile, icon_url, token, secret, website, location FROM users WHERE name = ?", name)
	err := row.Scan(&this.Id, &this.Name, &this.Profile, &this.IconURL, &this.Token, &this.Secret, &this.Website, &this.Location)

	if err != nil {
		return err
	}

	return nil
}

func (this *User) Update() error {

	_, err := DB.Exec("UPDATE users SET name = ?, profile = ?, icon_url = ?, token = ?, secret = ?, website = ?, location = ? WHERE id = ?",
		this.Name, this.Profile, this.IconURL, this.Token, this.Secret, this.Website, this.Location, this.Id)

	if err != nil {
		return err
	}

	return nil
}

func (this *User) Create() (int, error) {

	result, err := DB.Exec("INSERT INTO users ( name, profile, icon_url, token, secret, website, location) VALUES ( ?, ?, ?, ?, ?, ?, ? )", this.Name, this.Profile, this.IconURL, this.Token, this.Secret, this.Website, this.Location)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

// TODO: うまく動いてない？
func ExistsUser(name string) bool {

	_, err := DB.Exec("SELECT id FROM users WHERE name = ?", name)
	if err != nil {
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

func GetUserIdFromName(name string) (int, error) {

	var id int

	row := DB.QueryRow("SELECT id FROM users WHERE name = ?", name)
	err := row.Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
