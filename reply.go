package main

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const startcommand string = "/start"
const nextmatchcommand string = "/nextmatch"

func reply(update *tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	if update.Message == nil {
		return nil, errors.New("Not a message update")
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Text {
	case startcommand:
		{
			msg.Text = fmt.Sprintf(`ברוכים הבאים ❤️🖤❤️🖤!
		כדי לשאול אותי מתי המשחק הבא, שלחו לי את הפקודה %s`, nextmatchcommand)
		}
	case nextmatchcommand:
		{
			msg.Text = "המשחק הבא יהיה ב..."
		}
	default:
		{
			if update.Message.Chat.IsPrivate() {
				msg.Text = fmt.Sprintf(`מצטער, אני לא יודע מה לעשות עם ״%s״
			יש רק דבר אחד שאני יודע לעשות, אבל אני עושה אותו ממש טוב 😇
			כדי לראות אותי בפעולה, שלחו לי %s`, update.Message.Text, nextmatchcommand)
				msg.ReplyToMessageID = update.Message.MessageID
			} else {
				return nil, errors.New("Nothing to send")
			}
		}
	}

	return &msg, nil
}
