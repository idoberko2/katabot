package main

import (
	"fmt"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestReply_nil(t *testing.T) {
	u := tgbotapi.Update{
		UpdateID: 17,
	}

	_, err := reply(&u)
	if err == nil {
		t.Fatal("Expected nil, got ", err)
	}
}

func TestReply_start(t *testing.T) {
	starttext := fmt.Sprintf(`ברוכים הבאים ❤️🖤❤️🖤!
		כדי לשאול אותי מתי המשחק הבא, שלחו לי את הפקודה %s`, nextmatchcommand)
	u := tgbotapi.Update{
		UpdateID: 18,
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 19,
			},
			Text: startcommand,
		},
	}

	msg, err := reply(&u)
	if err != nil {
		t.Fatal("Expected err to be nil, got ", err)
	}

	if msg.ChatID != 19 {
		t.Fatal("Expected ChatID to be 19, got ", msg.ChatID)
	}

	if msg.Text != starttext {
		t.Fatalf(`Expected %s response to be
		"%s"
		got:
		"%s"`, startcommand, starttext, msg.Text)
	}
}

func TestReply_nextmatch(t *testing.T) {
	nextmatchtext := "המשחק הבא יהיה ב..."
	u := tgbotapi.Update{
		UpdateID: 19,
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 20,
			},
			Text: nextmatchcommand,
		},
	}

	msg, err := reply(&u)
	if err != nil {
		t.Fatal("Expected err to be nil, got ", err)
	}

	if msg.Text != nextmatchtext {
		t.Fatalf(`Expected %s response to be
		"%s"
		got:
		"%s"`, nextmatchcommand, nextmatchtext, msg.Text)
	}
}

func TestReply_default_group(t *testing.T) {
	u := tgbotapi.Update{
		UpdateID: 20,
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID:   21,
				Type: "group",
			},
			Text: "unknown commandddd",
		},
	}

	_, err := reply(&u)
	if err == nil {
		t.Fatal("Expected nil, got ", err)
	}
}

func TestReply_default_private(t *testing.T) {
	uc := "unknown commandddd"
	privatewelcome := fmt.Sprintf(`מצטער, אני לא יודע מה לעשות עם ״%s״
			יש רק דבר אחד שאני יודע לעשות, אבל אני עושה אותו ממש טוב 😇
			כדי לראות אותי בפעולה, שלחו לי %s`, uc, nextmatchcommand)
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

	msg, err := reply(&u)
	if err != nil {
		t.Fatal("Expected err to be nil, got ", err)
	}

	if msg.Text != privatewelcome {
		t.Fatalf(`Expected %s response to be
		"%s"
		got:
		"%s"`, nextmatchcommand, privatewelcome, msg.Text)
	}
}
