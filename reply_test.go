package main

import (
	"fmt"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type fakeBot struct {
	expectedMessage string
	testedCommand   string
	callCounter     int
	t               *testing.T
}

func (bot *fakeBot) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	bot.callCounter++

	switch v := c.(type) {
	case *tgbotapi.MessageConfig:
		{
			if v.Text != bot.expectedMessage {
				bot.t.Fatalf(`Expected %s response to be
		"%s"
		got:
		"%s"`, bot.expectedMessage, bot.expectedMessage, v.Text)
			}
		}
	}

	return tgbotapi.Message{}, nil
}

func TestReply_nil(t *testing.T) {
	bot := fakeBot{
		t: t,
	}
	u := tgbotapi.Update{
		UpdateID: 17,
	}

	reply(&bot, &u)
	if bot.callCounter > 0 {
		t.Fatalf("Expected Send to not be called, but it was called %d times", bot.callCounter)
	}
}

func TestReply_start(t *testing.T) {
	bot := fakeBot{
		testedCommand: startcommand,
		expectedMessage: fmt.Sprintf(`ברוכים הבאים ❤️🖤❤️🖤!
		כדי לשאול אותי מתי המשחק הבא, שלחו לי את הפקודה %s`, nextmatchcommand),
		t: t,
	}
	u := tgbotapi.Update{
		UpdateID: 18,
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 19,
			},
			Text: startcommand,
		},
	}

	reply(&bot, &u)
}

func TestReply_nextmatch(t *testing.T) {
	bot := fakeBot{
		testedCommand:   nextmatchcommand,
		expectedMessage: "המשחק הבא יהיה ב...",
		t:               t,
	}
	u := tgbotapi.Update{
		UpdateID: 19,
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 20,
			},
			Text: nextmatchcommand,
		},
	}

	reply(&bot, &u)
}

func TestReply_default_group(t *testing.T) {
	bot := fakeBot{
		t:             t,
		testedCommand: "unknown commandddd",
	}
	u := tgbotapi.Update{
		UpdateID: 20,
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID:   21,
				Type: "group",
			},
			Text: bot.testedCommand,
		},
	}

	reply(&bot, &u)
	if bot.callCounter > 0 {
		t.Fatalf("Expected Send to not be called, but it was called %d times", bot.callCounter)
	}
}

func TestReply_default_private(t *testing.T) {
	uc := "unknown commandddd"
	bot := fakeBot{
		t:             t,
		testedCommand: uc,
		expectedMessage: fmt.Sprintf(`מצטער, אני לא יודע מה לעשות עם ״%s״
			יש רק דבר אחד שאני יודע לעשות, אבל אני עושה אותו ממש טוב 😇
			כדי לראות אותי בפעולה, שלחו לי %s`, uc, nextmatchcommand),
	}

	u := tgbotapi.Update{
		UpdateID: 20,
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID:   21,
				Type: "private",
			},
			Text: uc,
		},
	}

	reply(&bot, &u)
}
