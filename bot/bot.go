package bot

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Bot is the bot that will be used to send all requests to the GroupMe API
type Bot struct {
	GroupID int64 `json:"group_id"`
	AccessToken string `json:"access_token"`
}

// InitBot stores the group_id and token needed to make requests for a certain group
func InitBot(configFile string) *Bot {
	file, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Panic(err)
	}

	var bot Bot

	json.Unmarshal(file, &bot)

	return &bot
}