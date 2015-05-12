package bot

import (
	"hsproom/config"
	"hsproom/utils/twitter"

	"github.com/mrjones/oauth"
)

var client *twitter.OAuthClient
var token  *oauth.AccessToken

func init() {

	var err error
	client, err = twitter.NewOAuthClient(config.TwitterBotKey, config.TwitterBotSecret)
	if err != nil {
		panic(err)
	}

	token = new(oauth.AccessToken)
	token.Token  = config.TwitterBotAccessKey
	token.Secret = config.TwitterBotAccessSecret
}

func Del() {
}

func UpdateTweet(message string) error {
	return client.UpdateTweet(token, message)
}
