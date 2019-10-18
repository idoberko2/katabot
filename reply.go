package main

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func reply(update *tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	if update.Message == nil {
		return nil, errors.New("Not a message update")
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if "הבא" == update.Message.Text || "מתי המשחק הבא?" == update.Message.Text || "next" == update.Message.Text {
		msg.Text = "המשחק הבא יהיה ב..."
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	} else {
		msg.Text = fmt.Sprintf(`מצטער, אני לא יודע מה לעשות עם ״%s״
		יש רק דבר אחד שאני יודע לעשות, אבל אני עושה אותו ממש טוב 😇`, update.Message.Text)
		msg.ReplyMarkup = numericKeyboard
	}

	msg.ReplyToMessageID = update.Message.MessageID

	return &msg, nil
}
