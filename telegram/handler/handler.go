package handler

import (
	"context"
	"log"

	bt "github.com/SakoDroid/telego"
	objs "github.com/SakoDroid/telego/objects"
	app "github.com/lordvidex/gomoney/telegram/application"
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
	b.bt.AddHandler("/createuser", b.CreateUser, "private")
	b.bt.AddHandler("/createuser", b.CreateGroupUser, "group", "supergroup")
	b.bt.AddHandler("/getuser", b.GetUser, "all")
	b.bt.AddHandler("/login", b.Login, "private")
	b.bt.AddHandler("/login", b.LoginGroup, "group", "supergroup")
	b.bt.AddHandler("/createaccount", b.CreateAccount, "all")
	b.bt.AddHandler("/getaccounts", b.GetAccounts, "all")
	b.bt.AddHandler("/deleteaccount", b.DeleteAccount, "all")
	b.bt.AddHandler("/transfer", b.Transfer, "all")
	b.bt.AddHandler("/deposit", b.Deposit, "all")
	b.bt.AddHandler("/withdraw", b.Withdraw, "all")
	b.bt.AddHandler("/gettransaction", b.GetTransaction, "all")
	b.bt.AddHandler("/getsummary", b.GetSummary, "all")
	b.bt.AddHandler("/help", b.HelpMessage, "all")
}

func (b *botHandler) HelpMessage(u *objs.Update) {
	b.simpleSend(u.Message.Chat.Id,
		`This is the GoMoni bot. Below are the commands you can use:
		/createuser - Create a new user account đ
		/login - Login to your account with your phone number đ
		/getuser - Get user details đ
		
		/createaccount - Create a new account đ
		/getaccounts - Get all your accounts â
		/deleteaccount - Delete an account đī¸
		/getsummary - Get transaction summary for all accounts đˇ
		
		/gettransaction - Get transaction details for a single account đˇ
		/transfer - Transfer money between accounts âī¸
		/deposit - Deposit money into an account âŦī¸
		/withdraw - Withdraw money from an account âŦī¸

		/help - Show this message âšī¸
		`,
		0,
	)
}

func (b *botHandler) simpleSend(chatID int, text string, replyTo int) {
	_, err := b.bt.SendMessage(chatID, text, "", replyTo, false, false)
	if err != nil {
		log.Println(err)
	}
}
