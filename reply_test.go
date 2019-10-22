package main

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type fakeBot struct {
	expectedMessage string
	testedCommand   string
	callCounter     int
	t               *testing.T
}

type fakeGamesFetcher struct {
	mock.Mock
}

func (gf *fakeGamesFetcher) GetNextKatamonGame(ctx context.Context) (*Round, *Game, error) {
	args := gf.Called(ctx)

	var arg0 *Round

	if args.Get(0) != nil {
		arg0 = args.Get(0).(*Round)
	}

	var arg1 *Game

	if args.Get(1) != nil {
		arg1 = args.Get(1).(*Game)
	}

	return arg0, arg1, args.Error(2)
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
	ctx := context.Background()
	bot := fakeBot{
		t: t,
	}
	u := tgbotapi.Update{
		UpdateID: 17,
	}

	reply(ctx, &bot, &u, &fakeGamesFetcher{})
	if bot.callCounter > 0 {
		t.Fatalf("Expected Send to not be called, but it was called %d times", bot.callCounter)
	}
}

func TestReply_start(t *testing.T) {
	ctx := context.Background()
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

	reply(ctx, &bot, &u, &fakeGamesFetcher{})
}

func TestReply_nextmatch(t *testing.T) {
	ctx := context.Background()
	fg := fakeGamesFetcher{}
	gd := time.Now().Add(time.Hour * 24)
	g := Game{
		Date:      gd,
		HomeTeam:  "הפועל קטמון י-ם",
		GuestTeam: "Beitar",
		Stadium:   "Teddy",
	}
	r := Round{
		LeagueID:    "45",
		SeasonID:    "21",
		RoundID:     "8",
		IsCompleted: false,
		Games:       []Game{g},
	}
	fg.On("GetNextKatamonGame", mock.Anything).Return(&r, &g, nil)
	bot := fakeBot{
		testedCommand: nextmatchcommand,
		expectedMessage: fmt.Sprintf(`המשחק הבא - מחזור %s
%s - %s
מיקום: %s
יום %s, %s, %s`, r.RoundID, g.HomeTeam, g.GuestTeam, g.Stadium, translateDay(g.Date.Format("Monday")), g.Date.Format("02/01"), g.Date.Format("15:04")),
		t: t,
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

	reply(ctx, &bot, &u, &fg)
}

func TestReply_nextmatch_error(t *testing.T) {
	ctx := context.Background()
	fg := fakeGamesFetcher{}
	fg.On("GetNextKatamonGame", mock.Anything).Return(nil, nil, errors.New("Some error occurred fetching next game"))
	bot := fakeBot{
		testedCommand: nextmatchcommand,
		expectedMessage: `משהו קרה ואני לא מצליח למצוא את המשחק הבא 🤔
נקווה שבפעם הבאה שתנסו אצליח אבל אין לדעת ¯\_(ツ)_/¯`,
		t: t,
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

	reply(ctx, &bot, &u, &fg)
}

func TestReply_default_group(t *testing.T) {
	ctx := context.Background()
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

	reply(ctx, &bot, &u, &fakeGamesFetcher{})
	if bot.callCounter > 0 {
		t.Fatalf("Expected Send to not be called, but it was called %d times", bot.callCounter)
	}
}

func TestReply_default_private(t *testing.T) {
	ctx := context.Background()
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

	reply(ctx, &bot, &u, &fakeGamesFetcher{})
}
