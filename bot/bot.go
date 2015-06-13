package bot

import (
	"strconv"
	"time"
	"os"

	"../config"
	"../utils/twitter"
	"../utils/log"
	"../models"

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

	go func() {
		var now time.Time
		var err error

		for {
			now = time.Now()

			// 12時0分
			if now.Hour() == 12 && now.Minute() == 0 {
				err = UpdateRankingForDay()
				log.InfoStr(os.Stdout, "Tweeted daily ranking.")
			}

			// 金曜日12時0分
			if now.Weekday() == time.Friday && now.Hour() == 12 && now.Minute() == 0 {
				err = UpdateRankingForWeek()
				log.InfoStr(os.Stdout, "Tweeted weekly ranking.")
			}

			// 1日12時0分
			if now.Day() == 1 && now.Hour() == 12 && now.Minute() == 0 {
				err = UpdateRankingForMonth()
				log.InfoStr(os.Stdout, "Tweeted monthly ranking.")
			}

			if err != nil {
				log.Fatal(os.Stdout, err)
			}

			time.Sleep(60 * time.Second)
		}
	}()
}

func Del() {
}

func UpdateTweet(message string) error {
	return client.UpdateTweet(token, message)
}

func UpdateRankingForWeek() error {
	var programs []models.ProgramInfo

	count, err := models.GetProgramRankingForWeek(&programs, 0, 3)
	if err != nil {
		return err
	}

	if count < 1 {
		return nil
	}

	message := "今週のプログラムランキング! #hsproom"

	for i := 0; i < count; i++ {
		message += "\n"
		message += strconv.Itoa(i+1)+"位: "
		message += programs[i].Title + " "
		message += "by " + programs[i].UserName + " "+ config.SiteURL +"/program/view/?p="+ strconv.Itoa(programs[i].Id)
	}

	return UpdateTweet(message)
}

func UpdateRankingForMonth() error {
	var programs []models.ProgramInfo

	count, err := models.GetProgramRankingForMonth(&programs, 0, 3)
	if err != nil {
		return err
	}

	if count < 1 {
		return nil
	}

	message := "今月のプログラムランキング! #hsproom"

	for i := 0; i < count; i++ {
		message += "\n"
		message += strconv.Itoa(i+1)+"位: "
		message += programs[i].Title + " "
		message += "by " + programs[i].UserName + " "+ config.SiteURL +"/program/view/?p="+ strconv.Itoa(programs[i].Id)
	}

	return UpdateTweet(message)
}

func UpdateRankingForDay() error {
	var programs []models.ProgramInfo

	count, err := models.GetProgramRankingForDay(&programs, 0, 3)
	if err != nil {
		return err
	}

	if count < 1 {
		return nil
	}

	message := "今日のプログラムランキング! #hsproom"

	for i := 0; i < count; i++ {
		message += "\n"
		message += strconv.Itoa(i+1)+"位: "
		message += programs[i].Title + " "
		message += "by " + programs[i].UserName + " "+ config.SiteURL +"/program/view/?p="+ strconv.Itoa(programs[i].Id)
	}

	return UpdateTweet(message)
}

