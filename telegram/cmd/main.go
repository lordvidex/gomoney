package main

import (
	"fmt"
	"time"

	bt "github.com/SakoDroid/telego"
	cfg "github.com/SakoDroid/telego/configs"
	objs "github.com/SakoDroid/telego/objects"
)

const (
	apiToken = "5817152401:AAGQKxskK9tN_I7oordUsblu1z5hQyzlmfI"
)

//gomoneyztest_bot

func main() {

	bot, err := bt.NewBot(cfg.Default(apiToken))

	if err == nil {
		err = bot.Run()
		if err == nil {
			start(bot)
		}
	}
}

func start(bot *bt.Bot) {

	//The general update channel.
	updateChannel := bot.GetUpdateChannel()
	//Adding a handler. Everytime the bot receives message "hi" in a private chat, it will respond "hi to you too".
	bot.AddHandler("hi", func(u *objs.Update) {
		if u.Message.Chat.Id == 672015206 {
			time.Sleep(time.Minute)
		}
		_, err := bot.SendMessage(u.Message.Chat.Id, "hi to you too", "", u.Message.MessageId, false, false)
		if err != nil {
			fmt.Println(err)
		}
	}, "private")

	//Monitores any other update. (Updates that don't contain text message "hi" in a private chat)
	for {
		update := <-*updateChannel
		fmt.Println("The chat id is", update.Message.Chat.Id)
		_, err := bot.SendMessage(update.Message.Chat.Id, "You said: "+update.Message.Text, "", update.Message.MessageId, false, false)
		if err != nil {
			fmt.Println(err)
		}
	}
}
