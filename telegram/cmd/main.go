package main

import (
	"context"
	"fmt"
	bt "github.com/SakoDroid/telego"
	cfg "github.com/SakoDroid/telego/configs"
	"github.com/lordvidex/gomoney/pkg/config"
	mgrpc "github.com/lordvidex/gomoney/telegram/adapters/grpc"
	"github.com/lordvidex/gomoney/telegram/application"
	"github.com/lordvidex/gomoney/telegram/handler"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const (
	apiToken = "5817152401:AAGQKxskK9tN_I7oordUsblu1z5hQyzlmfI"
)

//gomoneyztest_bot

func main() {
	// read configs
	c := config.New()

	bot, err := bt.NewBot(cfg.Default(apiToken))
	if err != nil {
		log.Fatal(err)
	}
	err = bot.Run()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create grpc service
	grpconn, err := connectGRPC(c)
	if err != nil {
		log.Fatal(err)
	}

	srv := mgrpc.New(grpconn)
	uc := application.New(srv)
	start(bot, uc, ctx)
}

func start(bot *bt.Bot, app *application.UseCases, ctx context.Context) {

	//The general update channel.
	updateChannel := bot.GetUpdateChannel()
	h := handler.NewBotHandler(bot, app, ctx)
	h.Register()
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

func connectGRPC(c *config.Config) (*grpc.ClientConn, error) {
	server := c.Get("GRPC_SERVER")
	if server == "" {
		return nil, errors.New("key 'GRPC_SERVER' not set")
	}
	return grpc.Dial(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
