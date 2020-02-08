package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	t := os.Getenv("TELEGRAM_APITOKEN")
	if t == "" {
		log.Fatal("Missing telegram API token")
	}

	bot, err := tgbotapi.NewBotAPI(t)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		// listen for SIGINT and SIGTERM events
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		bot.StopReceivingUpdates()
	}()

	gh := os.Getenv("GRAPHQL_HOST")
	if gh == "" {
		log.Fatal("Missing GraphQL host")
	}

	bot.Debug = os.Getenv("DEBUG") == "1"

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, err := bot.GetUpdatesChan(updateConfig)

	ctx := context.Background()
	gf := GamesFetcherInst{
		Client: APIClient{
			URL:  fmt.Sprintf("%s%s/graphql", gh, os.Getenv("GRAPHQL_PORT")),
			Mime: "application/json",
		},
	}

	s := botSender{
		bot: bot,
	}

	for update := range updates {
		reply(ctx, &s, &update, &gf)
	}
}
