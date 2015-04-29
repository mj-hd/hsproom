package google

import (
	"sync"
	"errors"
	"encoding/json"
	"encoding/binary"
	"crypto/rand"
	"io/ioutil"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func init() {
}

type OAuth2RequestTokenPool struct {
	tokens []string
	mutex  sync.Mutex
}

func NewOAuth2RequestTokenPool() *OAuth2RequestTokenPool {
	var tokenPool OAuth2RequestTokenPool

	tokenPool.tokens = make([]string, 0)

	return &tokenPool
}

func (this *OAuth2RequestTokenPool) Add(token string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.tokens = append(this.tokens, token)
}

func (this *OAuth2RequestTokenPool) Del(token string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	var i int
	var v string

	for i, v = range this.tokens {
		if v == token {
			break
		}
	}

	if v != token {
		return
	}

	this.tokens = append(this.tokens[:i], this.tokens[(i+1):]...)
}

func (this *OAuth2RequestTokenPool) Get(token string) string {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	for _, v := range this.tokens {
		if v == token {
			return v
		}
	}

	return ""
}


type OAuth2Client struct {
	Config *oauth2.Config
	tokenPool *OAuth2RequestTokenPool
}

func NewOAuth2Client(consumerKey string, consumerSecret string) (*OAuth2Client, error) {

	client := new(OAuth2Client)

	client.Config = &oauth2.Config{
		ClientID:     consumerKey,
		ClientSecret: consumerSecret,
		RedirectURL:  "",
		Scopes: []string{
			"https://www.googleapis.com/auth/plus.login",
			"https://www.googleapis.com/auth/plus.me",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	client.tokenPool = NewOAuth2RequestTokenPool()

	return client, nil
}

func (this *OAuth2Client) GetAuthURL(callbackUrl string) (string, error) {

	this.Config.RedirectURL = callbackUrl

	// ランダムな文字列を生成
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	verifier := strconv.FormatUint(n, 36)

	this.tokenPool.Add(verifier)

	url := this.Config.AuthCodeURL(verifier)

	return url, nil
}

func (this *OAuth2Client) GetToken(verifier string, code string) (*oauth2.Token, error) {

	token := this.tokenPool.Get(verifier)

	if token == "" {
		return nil, errors.New("リクエストトークンが不正。")
	}

	this.tokenPool.Del(verifier)

	return this.Config.Exchange(oauth2.NoContext, code)
}

func (this *OAuth2Client) GetUser(token *oauth2.Token) (*Userinfo, error) {
	client := this.Config.Client(oauth2.NoContext, token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()


	if response.StatusCode != 200 {
		return nil, errors.New("Unexpected HTTP Status Code " + strconv.Itoa(response.StatusCode) + " has returned!")
	}

	userinfo    := new(Userinfo)
	buffer, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buffer, userinfo)

	if err != nil {
		return nil, err
	}

	return userinfo, nil
}
