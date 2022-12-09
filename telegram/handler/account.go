package handler

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	objs "github.com/SakoDroid/telego/objects"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

func (b *botHandler) CreateAccount(u *objs.Update) {
	chatID := strconv.Itoa(u.Message.Chat.Id)
	ch, err := b.bt.AdvancedMode().RegisterChannel(chatID, "message")
	if err != nil {
		log.Println(err)
		return
	}
	defer b.bt.AdvancedMode().UnRegisterChannel(chatID, "message")
	_, err = b.bt.SendMessage(u.Message.Chat.Id, "1/3: Enter the title of the account. /cancel to terminate.", "", u.Message.MessageId, false, false)
	if err != nil {
		log.Println(err)
		return
	}
	u = <-*ch
	title := u.Message.Text
	if title == "/cancel" {
		b.simpleSend(u.Message.Chat.Id, "Account creation cancelled.", u.Message.MessageId)
		return
	}
	_, err = b.bt.SendMessage(u.Message.Chat.Id, "2/3: Enter the description of the account. /cancel to terminate", "", u.Message.MessageId, false, false)
	if err != nil {
		log.Println(err)
		return
	}
	u = <-*ch
	desc := u.Message.Text
	if title == "/cancel" {
		b.simpleSend(u.Message.Chat.Id, "Account creation cancelled.", u.Message.MessageId)
		return
	}

	// currency
	ik := b.bt.CreateKeyboard(true, true, false, "")
	ik.AddButton("USD", 1)
	ik.AddButton("RUB", 2)
	ik.AddButton("NGN", 2)
	ik.AddButton("Cancel", 3)
	b.bt.AdvancedMode().ASendMessage(u.Message.Chat.Id, "3/3: Choose the currency of the account.", "", u.Message.MessageId, false, false, nil, true, false, ik)
	u = <-*ch
	currency := u.Message.Text
	switch currency {
	case "USD", "RUB", "NGN":
	case "Cancel":
		b.simpleSend(u.Message.Chat.Id, "Account creation cancelled.", u.Message.MessageId)
		return
	default:
		b.simpleSend(u.Message.Chat.Id, "Invalid currency. Please try again.", u.Message.MessageId)
		return
	}
	curr := gomoney.Currency(currency)
	id, err := b.a.CreateAccount(b.ctx, strconv.Itoa(u.Message.Chat.Id), &gomoney.Account{
		Title:       title,
		Description: desc,
		Currency:    curr,
	})
	if err != nil {
		log.Println(err)
		b.simpleSend(u.Message.Chat.Id, "An error occured. Please try again.", u.Message.MessageId)
		return
	}
	b.simpleSend(u.Message.Chat.Id, fmt.Sprintf("Account created successfully. ID: %d", id), u.Message.MessageId)
}

func (b *botHandler) GetAccounts(u *objs.Update) {
	accs, err := b.a.GetAccounts(b.ctx, strconv.Itoa(u.Message.Chat.Id))
	if err != nil {
		b.simpleSend(u.Message.Chat.Id, "An error occurred while fetching accounts. Please try again.", u.Message.MessageId)
	}
	b.bt.SendMessage(u.Message.Chat.Id, BeautifulAccounts(accs), "MarkdownV2", u.Message.MessageId, false, false)
}

func (b *botHandler) DeleteAccount(u *objs.Update) {
	chatID := strconv.Itoa(u.Message.Chat.Id)
	ch, err := b.bt.AdvancedMode().RegisterChannel(chatID, "message")
	if err != nil {
		log.Println(err)
		return
	}
	defer b.bt.AdvancedMode().UnRegisterChannel(chatID, "message")

	accs, err := b.a.GetAccounts(b.ctx, chatID)
	if err != nil {
		log.Println(err)
		return
	}
	if len(accs) == 0 {
		b.simpleSend(u.Message.Chat.Id, "You have no accounts to delete.", u.Message.MessageId)
		return
	}
	kb := b.bt.CreateKeyboard(true, true, false, "1...")
	for i, acc := range accs {
		kb.AddButton(fmt.Sprintf("%d. %s", i+1, acc.Title), (i/2)+1)
	}
	kb.AddButton("Cancel", (len(accs)/2)+1)
	_, err = b.bt.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Click the account you want to delete. Cancel to terminate", "", u.Message.MessageId, false, false, nil, true, false, kb)
	if err != nil {
		log.Println(err)
		return
	}
	u = <-*ch
	res := u.Message.Text
	if res == "Cancel" {
		b.simpleSend(u.Message.Chat.Id, "Account deletion cancelled.", u.Message.MessageId)
		return
	}
	idx, err := strconv.Atoi(strings.Split(res, ".")[0])
	if err != nil {
		log.Println(err)
		return
	}
	err = b.a.DeleteAccount(b.ctx, accs[idx-1].Id, chatID)
	if err != nil {
		log.Println(err)
		b.simpleSend(u.Message.Chat.Id, "An error occurred while deleting account. Please try again.", u.Message.MessageId)
		return
	}
	b.simpleSend(u.Message.Chat.Id, "Account deleted successfully.", u.Message.MessageId)
}
