package main

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const startcommand string = "/start"
const nextmatchcommand string = "/nextmatch"

type sender interface {
	SendText(cid int64, t string) (tgbotapi.Message, error)
}

type botSender struct {
	bot *tgbotapi.BotAPI
}

func (s botSender) SendText(cid int64, t string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(cid, t)
	return s.bot.Send(&msg)
}

func translateDay(d string) string {
	switch d {
	case "Sunday":
		return "ראשון"
	case "Monday":
		return "שני"
	case "Tuesday":
		return "שלישי"
	case "Wednesday":
		return "רביעי"
	case "Thursday":
		return "חמישי"
	case "Friday":
		return "שישי"
	case "Saturday":
		return "שבת"
	}

	return ""
}

func reply(ctx context.Context, bot sender, update *tgbotapi.Update, gf GamesFetcher) {
	if update.Message == nil {
		log.Println("Not a message update")
		return
	}

	msg := ""

	switch update.Message.Text {
	case startcommand:
		{
			msg = fmt.Sprintf(`ברוכים הבאים ❤️🖤❤️🖤!
כדי לשאול אותי מתי המשחק הבא, שלחו לי את הפקודה %s`, nextmatchcommand)
		}
	case nextmatchcommand:
		{
			r, g, err := gf.GetNextKatamonGame(ctx)
			if err != nil {
				msg = `משהו קרה ואני לא מצליח למצוא את המשחק הבא 🤔
נקווה שבפעם הבאה שתנסו אצליח אבל אין לדעת ¯\_(ツ)_/¯`
			} else {
				msg = fmt.Sprintf(`המשחק הבא - מחזור %s
%s - %s
מיקום: %s
יום %s, %s, %s`, r.RoundID, g.HomeTeam, g.GuestTeam, g.Stadium, translateDay(g.Date.Format("Monday")), g.Date.Format("02/01"), g.Date.Format("15:04"))
			}
		}
	default:
		{
			if update.Message.Chat.IsPrivate() {
				msg = fmt.Sprintf(`מצטער, אני לא יודע מה לעשות עם ״%s״
יש רק דבר אחד שאני יודע לעשות, אבל אני עושה אותו ממש טוב 😇
כדי לראות אותי בפעולה, שלחו לי %s`, update.Message.Text, nextmatchcommand)
				// msg.ReplyToMessageID = update.Message.MessageID
			} else {
				log.Println("Nothing to send")
				return
			}
		}
	}

	if _, err := bot.SendText(update.Message.Chat.ID, msg); err != nil {
		log.Fatal(err)
	}
}
