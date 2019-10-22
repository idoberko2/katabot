package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

type sender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

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
			URL:  fmt.Sprintf("http://%s%s/graphql", gh, os.Getenv("GRAPHQL_PORT")),
			Mime: "application/json",
		},
	}

	for update := range updates {
		go reply(ctx, bot, &update, &gf)
	}
}
