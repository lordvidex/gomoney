package handler

import (
	"context"
	bt "github.com/SakoDroid/telego"
	app "github.com/lordvidex/gomoney/telegram/application"
	"log"
)

type botHandler struct {
	bt  *bt.Bot
	a   *app.UseCases
	ctx context.Context
}

func NewBotHandler(bt *bt.Bot, a *app.UseCases, ctx context.Context) *botHandler {
	return &botHandler{
		bt:  bt,
		a:   a,
		ctx: ctx,
	}
}

func (b *botHandler) Register() {
	b.bt.AddHandler("/createAccount", b.CreateAccount, "private")
}

func (b *botHandler) simpleSend(chatID int, text string, replyTo int) {
	_, err := b.bt.SendMessage(chatID, text, "", replyTo, false, false)
	if err != nil {
		log.Println(err)
	}
}
