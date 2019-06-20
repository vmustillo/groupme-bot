package bot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"fmt"
)

// Bot is the bot that will be used to send all requests to the GroupMe API
type Bot struct {
	GroupID int64 `json:"group_id"`
	AccessToken string `json:"access_token"`
}

// InitBot stores the group_id and token needed to make requests for a certain group
func InitBot(configFile string) *Bot {
	file, err := ioutil.ReadFile(configFile)
	Handle(err)

	var bot Bot
	json.Unmarshal(file, &bot)

	return &bot
}

// GetMessages makes a get request to the GroupMe API and returns a specified number of messages
func (bot *Bot) GetMessages(numMsgs int) []byte {
	queryString := fmt.Sprintf("https://api.groupme.com/v3/groups/%d/messages?limit=%d&token=%s", bot.GroupID, numMsgs, bot.AccessToken)
	
	res, err := http.Get(queryString)
	Handle(err)

	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return data
}

// Handle is a simple error handler to log errors call a panic
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
