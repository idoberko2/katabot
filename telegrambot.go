package main

import (
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
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = os.Getenv("DEBUG") == "1"

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, err := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		go reply(bot, &update)
	}
}
