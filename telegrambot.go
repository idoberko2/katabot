package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("מתי המשחק הבא?"),
	),
)

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
		msg, err := reply(&update)

		if err != nil {
			log.Println(err)
			continue
		}

		if _, err := bot.Send(msg); err != nil {
			log.Fatal(err)
		}
	}
}
