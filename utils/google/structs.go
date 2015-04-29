package google

type Userinfo struct {
	IdString string `json:"id"`
	Name string `json:"name"`
	GivenName string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Link string `json:"link"`
	Picture string `json:"picture"`
	Locale string `json:"locale"`
}
