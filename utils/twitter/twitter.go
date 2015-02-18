package twitter

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/mrjones/oauth"
)

func init() {
}

type Client struct {
	httpClient     *http.Client
	consumerKey    string
	consumerSecret string
	accessToken    string
}

func NewClient(consumerKey string, consumerSecret string) (*Client, error) {

	var client Client

	client.consumerKey = consumerKey
	client.consumerSecret = consumerSecret

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}

	client.httpClient = &http.Client{
		Transport: tr,
	}

	return &client, client.validate()
}

type oauth2TokenMember struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

func (this *Client) validate() error {
	bearerToken := base64.StdEncoding.EncodeToString([]byte(this.consumerKey + ":" + this.consumerSecret))

	request, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return err
	}

	request.Header.Add("User-Agent", "HSPRoom")
	request.Header.Add("Authorization", "Basic "+bearerToken)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Cache-Control", "no-cache")

	response, err := this.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {

		return errors.New("HTTP Status " + response.Status + " has returned.")
	}

	var result oauth2TokenMember
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)

	if err != nil {
		return err
	}

	if result.TokenType != "bearer" {
		return errors.New("Unknown token type has returned " + result.TokenType)
	}

	this.accessToken = result.AccessToken

	return nil
}

func (this *Client) Close() {

	bearerToken := base64.StdEncoding.EncodeToString([]byte(this.consumerKey + ":" + this.consumerSecret))

	request, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/invalidate_token", strings.NewReader("access_token="+this.accessToken))
	if err != nil {
		return
	}

	request.Header.Add("User-Agent", "HSPRoom")
	request.Header.Add("Authorization", "Basic "+bearerToken)
	request.Header.Add("Aeccept", "*/*")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := this.httpClient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return
	}

	this.accessToken = ""
}

func (this *Client) SearchTweets(query string, number int, offsetId int) (SearchResponse, error) {
	if this.accessToken == "" {
		return SearchResponse{}, errors.New("Must be initialized before calling this method.")
	}

	encodedQuery := url.QueryEscape(query)

	request, err := http.NewRequest("GET", "https://api.twitter.com/1.1/search/tweets.json?q="+encodedQuery+"&result_type=recent&count="+strconv.Itoa(number)+"&since_id="+strconv.Itoa(offsetId), nil)
	if err != nil {
		return SearchResponse{}, err
	}

	request.Header.Set("User-Agent", "HSPRoom")
	request.Header.Set("Authorization", "Bearer "+this.accessToken)

	response, err := this.httpClient.Do(request)
	if err != nil {
		return SearchResponse{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return SearchResponse{}, errors.New("HTTP Status " + response.Status + " has returned.")
	}

	var tweets SearchResponse

	err = json.NewDecoder(response.Body).Decode(&tweets)
	if err != nil {
		return SearchResponse{}, err
	}

	return tweets, nil
}

type requestTokenPool struct {
	tokens map[string]*oauth.RequestToken
	mutex  sync.Mutex
}

func newRequestTokenPool() *requestTokenPool {
	var tokenPool requestTokenPool

	tokenPool.tokens = make(map[string]*oauth.RequestToken)

	return &tokenPool
}

func (this *requestTokenPool) Add(token *oauth.RequestToken) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.tokens[token.Token] = token
}

func (this *requestTokenPool) Del(token string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	delete(this.tokens, token)
}

func (this *requestTokenPool) Get(token string) *oauth.RequestToken {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.tokens[token]
}

type OAuthClient struct {
	Config    *oauth.Consumer
	tokenPool *requestTokenPool
}

func NewOAuthClient(consumerKey string, consumerSecret string) (*OAuthClient, error) {
	client := new(OAuthClient)

	client.Config = oauth.NewConsumer(
		consumerKey,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		},
	)

	client.tokenPool = newRequestTokenPool()

	return client, nil
}

func (this *OAuthClient) GetAuthURL(callbackUrl string) (string, error) {

	requestToken, url, err := this.Config.GetRequestTokenAndUrl(callbackUrl)
	if err != nil {
		return "", err
	}

	this.tokenPool.Add(requestToken)

	return url, nil
}

func (this *OAuthClient) GetAccessToken(verifier string, token string) (*oauth.AccessToken, error) {
	requestToken := this.tokenPool.Get(token)

	if requestToken == nil {
		return nil, errors.New("リクエストトークンプールに値が見つからなかった。")
	}

	accessToken, err := this.Config.AuthorizeToken(requestToken, verifier)
	if err != nil {
		return &oauth.AccessToken{}, err
	}

	return accessToken, nil
}

func (this *OAuthClient) CheckUserCredentialsAndGetUser(accessToken *oauth.AccessToken) (*User, error) {

	response, err := this.Config.Get(
		"https://api.twitter.com/1.1/account/verify_credentials.json",
		map[string]string{},
		accessToken,
	)

	if err != nil {
		return &User{}, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return &User{}, errors.New("Unexpected HTTP Status Code " + strconv.Itoa(response.StatusCode) + " has returned!")
	}

	var user User
	buffer, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return &User{}, err
	}

	err = json.Unmarshal(buffer, &user)

	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

/*func main() {
	client, err := NewClient(config.TwitterKey, config.TwitterSecret)
	if err != nil {
		panic(err)
	}

	tweets, err := client.SearchTweets("HSP", 5, 1)
	if err != nil {
		panic(err)
	}

	for tweet := range tweets {
		js, _ := json.Marshal(tweet)
		println(string(js))
	}
}*/
