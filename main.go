package main

import (
	"github.com/vmustillo/groupme-response-bot/bot"
)

func main() {
	b := bot.InitBot("config.json")
	msgDump := b.GetMessages(20)
	res := bot.ParseMessages(msgDump)

	hotTakes := res.SearchMessages("hot take")
}
