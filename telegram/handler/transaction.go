package handler

import (
	"fmt"
	objs "github.com/SakoDroid/telego/objects"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	"github.com/lordvidex/gomoney/telegram/application"
	"log"
	"strconv"
	"strings"
)

func (b *botHandler) Withdraw(u *objs.Update) {
	chatID := strconv.Itoa(u.Message.Chat.Id)
	ch, err := b.bt.AdvancedMode().RegisterChannel(chatID, "message")
	if err != nil {
		log.Println(err)
		return
	}
	defer b.bt.AdvancedMode().UnRegisterChannel(chatID, "message")

	// get user accounts
	accs, err := b.a.GetAccounts(b.ctx, chatID)
	if err != nil {
		log.Println(err)
		return
	}
	if len(accs) == 0 {
		b.simpleSend(u.Message.Chat.Id, "You have no accounts to withdraw from.", u.Message.MessageId)
		return
	}
	// get account
	ik := b.bt.CreateKeyboard(true, true, false, "")
	for i, acc := range accs {
		ik.AddButton(fmt.Sprintf("%d. %s", i+1, acc.Title), (i/2)+1)
	}
	ik.AddButton("Cancel", (len(accs)/2)+1)
	_, err = b.bt.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Click the account you want to withdraw from. Cancel to terminate", "", u.Message.MessageId, false, false, nil, true, false, ik)
	if err != nil {
		log.Println(err)
		return
	}
	u = <-*ch
	res := u.Message.Text
	if res == "Cancel" {
		b.simpleSend(u.Message.Chat.Id, "Withdraw cancelled.", u.Message.MessageId)
		return
	}
	idx, err := strconv.Atoi(strings.Split(res, ".")[0])
	if err != nil {
		log.Println(err)
		return
	}
	// get amount
	b.simpleSend(u.Message.Chat.Id, "Enter the amount you want to withdraw.", u.Message.MessageId)
	u = <-*ch
	amount, err := strconv.ParseFloat(u.Message.Text, 64)
	if err != nil {
		b.simpleSend(u.Message.Chat.Id, "Invalid amount.", u.Message.MessageId)
		return
	}

	// withdraw money
	err = b.a.Withdraw(b.ctx, amount, accs[idx-1].Id, chatID)
	if err != nil {
		log.Println(err)
		b.simpleSend(u.Message.Chat.Id, "Withdraw failed.", u.Message.MessageId)
		return
	}
	b.simpleSend(u.Message.Chat.Id, "Withdrawal successful.", u.Message.MessageId)
}

func (b *botHandler) Transfer(u *objs.Update) {
	chatID := strconv.Itoa(u.Message.Chat.Id)
	ch, err := b.bt.AdvancedMode().RegisterChannel(chatID, "message")
	if err != nil {
		log.Println(err)
		return
	}
	defer b.bt.AdvancedMode().UnRegisterChannel(chatID, "message")

	// get user accounts
	accs, err := b.a.GetAccounts(b.ctx, chatID)
	if err != nil {
		log.Println(err)
		return
	}
	if len(accs) == 0 {
		b.simpleSend(u.Message.Chat.Id, "You have no accounts to transfer from.", u.Message.MessageId)
		return
	}
	// get account
	ik := b.bt.CreateKeyboard(true, true, false, "")
	for i, acc := range accs {
		ik.AddButton(fmt.Sprintf("%d. %s", i+1, acc.Title), (i/2)+1)
	}
	ik.AddButton("Cancel", (len(accs)/2)+1)

	// account from
	_, err = b.bt.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Click the account you want to transfer from. Cancel to terminate", "", u.Message.MessageId, false, false, nil, true, false, ik)
	if err != nil {
		log.Println(err)
		return
	}
	u = <-*ch
	res := u.Message.Text
	if res == "Cancel" {
		b.simpleSend(u.Message.Chat.Id, "Transfer cancelled.", u.Message.MessageId)
		return
	}
	idxFrom, err := strconv.Atoi(strings.Split(res, ".")[0])
	if err != nil {
		log.Println(err)
		return
	}

	// get account to
	_, err = b.bt.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Click the account you want to transfer to. Cancel to terminate", "", u.Message.MessageId, false, false, nil, true, false, ik)
	if err != nil {
		log.Println(err)
		return
	}
	u = <-*ch
	res = u.Message.Text
	if res == "Cancel" {
		b.simpleSend(u.Message.Chat.Id, "Transfer cancelled.", u.Message.MessageId)
		return
	}
	idxTo, err := strconv.Atoi(strings.Split(res, ".")[0])
	if err != nil {
		log.Println(err)
		return
	}

	// get amount
	b.simpleSend(u.Message.Chat.Id, "Enter the amount you want to transfer.", u.Message.MessageId)
	u = <-*ch
	amount, err := strconv.ParseFloat(u.Message.Text, 64)
	if err != nil {
		b.simpleSend(u.Message.Chat.Id, "Invalid amount.", u.Message.MessageId)
		return
	}

	// transfer money
	err = b.a.Transfer(
		b.ctx,
		application.TransferParam{
			From:   accs[idxFrom-1].Id,
			To:     accs[idxTo-1].Id,
			Amount: amount,
		},
		chatID,
	)
	if err != nil {
		log.Println(err)
		b.simpleSend(u.Message.Chat.Id, "Transfer failed.", u.Message.MessageId)
		return
	}
	b.simpleSend(u.Message.Chat.Id, "Transfer successful.", u.Message.MessageId)
}

func (b *botHandler) Deposit(u *objs.Update) {
	chatID := strconv.Itoa(u.Message.Chat.Id)
	ch, err := b.bt.AdvancedMode().RegisterChannel(chatID, "message")
	if err != nil {
		log.Println(err)
		return
	}
	defer b.bt.AdvancedMode().UnRegisterChannel(chatID, "message")

	// get user accounts
	accs, err := b.a.GetAccounts(b.ctx, chatID)
	if err != nil {
		log.Println(err)
		return
	}
	if len(accs) == 0 {
		b.simpleSend(u.Message.Chat.Id, "You have no accounts to deposit into.", u.Message.MessageId)
		return
	}
	// get account
	ik := b.bt.CreateKeyboard(true, true, false, "")
	for i, acc := range accs {
		ik.AddButton(fmt.Sprintf("%d. %s", i+1, acc.Title), (i/2)+1)
	}
	ik.AddButton("Cancel", (len(accs)/2)+1)
	_, err = b.bt.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Click the account you want to deposit into. Cancel to terminate", "", u.Message.MessageId, false, false, nil, true, false, ik)
	if err != nil {
		log.Println(err)
		return
	}
	u = <-*ch
	res := u.Message.Text
	if res == "Cancel" {
		b.simpleSend(u.Message.Chat.Id, "Deposit cancelled.", u.Message.MessageId)
		return
	}
	idx, err := strconv.Atoi(strings.Split(res, ".")[0])
	if err != nil {
		log.Println(err)
		return
	}
	// get amount
	b.simpleSend(u.Message.Chat.Id, "Enter the amount you want to deposit.", u.Message.MessageId)
	u = <-*ch
	amount, err := strconv.ParseFloat(u.Message.Text, 64)
	if err != nil {
		b.simpleSend(u.Message.Chat.Id, "Invalid amount.", u.Message.MessageId)
		return
	}

	// deposit money
	err = b.a.Deposit(b.ctx, amount, accs[idx-1].Id, chatID)
	if err != nil {
		log.Println(err)
		b.simpleSend(u.Message.Chat.Id, "Deposit failed.", u.Message.MessageId)
		return
	}
	b.simpleSend(u.Message.Chat.Id, "Deposit successful.", u.Message.MessageId)
}

func (b *botHandler) GetTransaction(u *objs.Update) {
	chatID := strconv.Itoa(u.Message.Chat.Id)
	ch, err := b.bt.AdvancedMode().RegisterChannel(chatID, "message")
	if err != nil {
		log.Println(err)
		return
	}
	defer b.bt.AdvancedMode().UnRegisterChannel(chatID, "message")

	// get the account to print transactions
	accs, err := b.a.GetAccounts(b.ctx, chatID)
	if err != nil {
		log.Println(err)
		return
	}
	if len(accs) == 0 {
		b.simpleSend(u.Message.Chat.Id, "You have no accounts.", u.Message.MessageId)
		return
	}
	ik := b.bt.CreateKeyboard(true, true, false, "")
	for i, acc := range accs {
		ik.AddButton(fmt.Sprintf("%d. %s", i+1, acc.Title), (i/2)+1)
	}
	ik.AddButton("Cancel", (len(accs)/2)+1)
	_, err = b.bt.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Click the account you want to view transactions from. Cancel to terminate", "", u.Message.MessageId, false, false, nil, true, false, ik)
	if err != nil {
		log.Println(err)
		return
	}
	u = <-*ch
	res := u.Message.Text
	if res == "Cancel" {
		b.simpleSend(u.Message.Chat.Id, "Action cancelled.", u.Message.MessageId)
		return
	}
	idx, err := strconv.Atoi(strings.Split(res, ".")[0])
	if err != nil {
		log.Println(err)
		return
	}
	acc := accs[idx-1]
	tx, err := b.a.GetAccountTransfers(b.ctx, acc.Id, chatID)
	if err != nil {
		log.Println(err)
		b.simpleSend(u.Message.Chat.Id, "Failed to get transactions.", u.Message.MessageId)
		return
	}
	b.bt.SendMessage(u.Message.Chat.Id, BeautifulTransactions(tx), "MarkdownV2", u.Message.MessageId, false, false)
}

func (b *botHandler) GetSummary(u *objs.Update) {
	chatID := strconv.Itoa(u.Message.Chat.Id)
	accs, err := b.a.GetAccounts(b.ctx, chatID)
	if err != nil {
		log.Println(err)
		return
	}
	if len(accs) == 0 {
		b.simpleSend(u.Message.Chat.Id, "You have no accounts.", u.Message.MessageId)
		return
	}
	m := make(map[int64]*gomoney.Account)
	for i := 0; i < len(accs); i++ {
		m[accs[i].Id] = &accs[i]
	}
	summaries, err := b.a.GetTransferSummary(b.ctx, chatID)
	if err != nil {
		log.Println(err)
		return
	}
	for i := 0; i < len(summaries); i++ {
		id := summaries[i].Account.Id
		summaries[i].Account = m[id]
	}
	b.bt.SendMessage(u.Message.Chat.Id, BeautifulTransferSummary(summaries), "MarkdownV2", u.Message.MessageId, false, false)
}
