package main

import (
	"fmt"

	"github.com/vmustillo/hot-take-bot/bot"
)

func main() {
	b := bot.InitBot("config.json")
	msgDump := b.GetMessages(20)
	res := bot.ParseMessages(msgDump)
	for _, v := range res.Data.Messages {
		fmt.Println(v.Message)
	}
}
