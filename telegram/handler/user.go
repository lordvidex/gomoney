package handler

import (
	"context"
	"fmt"
	"strings"

	objs "github.com/SakoDroid/telego/objects"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	app "github.com/lordvidex/gomoney/telegram/application"

	"log"
	"strconv"
	"time"
)

func (b *botHandler) CreateUser(u *objs.Update) {
	// check if the user exists OR create a new user
	// get phone number and pass keyboard to get contact
	k := b.bt.CreateKeyboard(true, true, false, "")
	k.AddContactButton("Send Contact", 1)
	_, err := b.bt.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Click the <b>Send Contact</b> button to send your phone number", "HTML", u.Message.MessageId, false, false, nil, true, true, k)
	if err != nil {
		log.Println(err)
		return
	}
	// create a new context from parent context
	ctx, cancel := context.WithTimeout(b.ctx, time.Minute)
	defer cancel()

	// register a new channel to listen for user replies
	ch, err := b.bt.AdvancedMode().RegisterChannel(strconv.Itoa(u.Message.Chat.Id), "message")
	defer b.bt.AdvancedMode().UnRegisterChannel(strconv.Itoa(u.Message.Chat.Id), "message") // close
	if err != nil {
		log.Println(err)
		return
	}

	select {
	case <-ctx.Done(): // if time expires or context was cancelled
		b.bt.SendMessage(u.Message.Chat.Id, "You took too long to respond", "", u.Message.MessageId, false, false)
		return
	case update := <-*ch: // if message is received from the channel

		if update.Message.Contact != nil {
			phone := update.Message.Contact.PhoneNumber
			if !strings.HasPrefix(phone, "+") {
				phone = "+" + phone
			}
			// create user
			id, err := b.a.CreateUser(b.ctx, app.CreateUserParam{
				Phone: phone,
				Name:  u.Message.From.FirstName + " " + u.Message.From.Lastname,
			})
			if err != nil {
				log.Println(err)
				if gomoney.ErrAlreadyExists.Is(err) {
					b.bt.SendMessage(u.Message.Chat.Id, "You already have an account with this phone number.", "", u.Message.MessageId, false, false)
					return
				}
				b.simpleSend(u.Message.Chat.Id, "An error occurred", u.Message.MessageId)
				return
			}
			_, err = b.bt.SendMessage(u.Message.Chat.Id, fmt.Sprintf("You have successfully registered with id %s.", id), "", u.Message.MessageId, false, false)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (b *botHandler) GetUser(u *objs.Update) {
	user, err := b.a.GetUserByChatID(b.ctx, strconv.Itoa(u.Message.Chat.Id))
	if err != nil {
		if gomoney.ErrNotFound.Is(err) {
			b.bt.SendMessage(u.Message.Chat.Id, `You don't have an account yet. Use /createuser to create one.`, "", u.Message.MessageId, false, false)
		}
		return
	}
	b.bt.SendMessage(u.Message.Chat.Id, BeautifulUserData(user), "MarkdownV2", u.Message.MessageId, false, false)
}

func (b *botHandler) Login(u *objs.Update) {
	// check if the user exists OR create a new user
	// get phone number and pass keyboard to get contact
	k := b.bt.CreateKeyboard(true, true, false, "")
	k.AddContactButton("Send Contact", 1)
	_, err := b.bt.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Click the <b>Send Contact</b> button to send your phone number", "HTML", u.Message.MessageId, false, false, nil, true, true, k)
	if err != nil {
		log.Println(err)
		return
	}
	// create a new context from parent context
	ctx, cancel := context.WithTimeout(b.ctx, time.Minute)
	defer cancel()

	// register a new channel to listen for user replies
	ch, err := b.bt.AdvancedMode().RegisterChannel(strconv.Itoa(u.Message.Chat.Id), "message")
	defer b.bt.AdvancedMode().UnRegisterChannel(strconv.Itoa(u.Message.Chat.Id), "message") // close
	if err != nil {
		log.Println(err)
		return
	}

	select {
	case <-ctx.Done(): // if time expires or context was cancelled
		b.bt.SendMessage(u.Message.Chat.Id, "You took too long to respond", "", u.Message.MessageId, false, false)
		return
	case update := <-*ch: // if message is received from the channel
		if update.Message.Contact != nil {
			phone := update.Message.Contact.PhoneNumber
			if !strings.HasPrefix(phone, "+") {
				phone = "+" + phone
			}
			// create user
			user, err := b.a.GetUserByPhone(b.ctx, phone, strconv.Itoa(u.Message.Chat.Id))
			if err != nil {
				log.Println(err)
				if gomoney.ErrNotFound.Is(err) {
					b.bt.SendMessage(u.Message.Chat.Id, "User with this phone number does not exist.", "", u.Message.MessageId, false, false)
					return
				}
				b.simpleSend(u.Message.Chat.Id, "An error occurred", u.Message.MessageId)
				return
			}
			_, err = b.bt.SendMessage(u.Message.Chat.Id, fmt.Sprintf("You have successfully logged in. Your id is %s.", user.ID), "", u.Message.MessageId, false, false)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
