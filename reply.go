package main

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const startcommand string = "/start"
const nextmatchcommand string = "/nextmatch"

func reply(ctx context.Context, bot sender, update *tgbotapi.Update, gf GamesFetcher) {
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
			_, g, _ := gf.GetNextKatamonGame(ctx)
			fmt.Printf("%+v\n", g)
			msg.Text = fmt.Sprintf(`המשחק הבא:
		%s - %s,
		מיקום: %s,
		זמן: %s`, g.HomeTeam, g.GuestTeam, g.Stadium, g.Date.Format(time.RFC3339))
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
