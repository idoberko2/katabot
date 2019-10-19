package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const startcommand string = "/start"
const nextmatchcommand string = "/nextmatch"

func reply(bot sender, update *tgbotapi.Update) {
	if update.Message == nil {
		log.Println("Not a message update")
		return
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
				log.Println("Nothing to send")
				return
			}
		}
	}

	if _, err := bot.Send(&msg); err != nil {
		log.Fatal(err)
	}
}
