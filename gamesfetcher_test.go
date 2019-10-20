package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fakeAPIClient struct {
	mock.Mock
}

func (fq *fakeAPIClient) Request(b map[string]interface{}, j interface{}) error {
	args := fq.Called(b, j)

	return args.Error(0)
}

func TestGetNextKatamonGame_error(t *testing.T) {
	fc := fakeAPIClient{}
	ctx := context.Background()
	fc.On("Request", mock.Anything, mock.Anything).Return(errors.New("Some error occurred"))

	_, _, err := GetNextKatamonGame(ctx, &fc)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGetNextKatamonGame_simple(t *testing.T) {
	fc := fakeAPIClient{}
	ctx := context.Background()
	gt := time.Now().Add(time.Hour * 12)
	fc.On("Request", map[string]interface{}{
		"operationName": "getNextRound",
		"query":         nextRoundQuery,
		"variables": map[string]int64{
			"timestamp": time.Now().Unix(),
		},
	}, &nextRoundResponse{}).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*nextRoundResponse)
		arg.Data = struct {
			NextRound roundResponse `json:"nextRound"`
		}{
			NextRound: roundResponse{
				League:      "45",
				Season:      "21",
				Round:       "7",
				IsCompleted: false,
				Games: []gameResponse{
					{
						HomeTeam:  "הפועל קטמון י-ם",
						GuestTeam: "Beitar",
						Stadium:   "Teddy",
						Date:      gt.Format(time.RFC3339),
					},
				},
			},
		}
	})

	r, g, err := GetNextKatamonGame(ctx, &fc)

	if err != nil {
		t.Fatal("Expected nil, got: ", err)
	}

	assert.Equal(t, r.LeagueID, "45", "Expected league 45")
	assert.Equal(t, r.SeasonID, "21", "Expected season 21")
	assert.Equal(t, r.RoundID, "7", "Expected round 7")
	assert.False(t, r.IsCompleted, "Expected false")
	assert.Equal(t, g.HomeTeam, "הפועל קטמון י-ם", "Expected Katamon to be the home team")
	assert.Equal(t, g.GuestTeam, "Beitar", "Expected Beitar to be the guest team")
	assert.Equal(t, g.Stadium, "Teddy", "Expected Teddy to be the stadium")
	assert.Equal(t, g.Date.Format(time.RFC3339), gt.Format(time.RFC3339), "Expected game time to be correct")
}

func TestGetNextKatamonGame_past(t *testing.T) {
	fc := fakeAPIClient{}
	ctx := context.Background()
	gt := time.Now().Add(time.Hour * -1)

	fc.On("Request", map[string]interface{}{
		"operationName": "getNextRound",
		"query":         nextRoundQuery,
		"variables": map[string]int64{
			"timestamp": time.Now().Unix(),
		},
	}, &nextRoundResponse{}).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*nextRoundResponse)
		arg.Data = struct {
			NextRound roundResponse `json:"nextRound"`
		}{
			NextRound: roundResponse{
				League:      "45",
				Season:      "21",
				Round:       "7",
				IsCompleted: false,
				Games: []gameResponse{
					{
						HomeTeam:  "הפועל קטמון י-ם",
						GuestTeam: "Beitar",
						Stadium:   "Teddy",
						Date:      gt.Format(time.RFC3339),
					},
				},
			},
		}
	})

	fc.On("Request", map[string]interface{}{
		"operationName": "getRound",
		"query":         getRoundQuery,
		"variables": map[string]string{
			"league": "45",
			"season": "21",
			"round":  "8",
		},
	}, &getRoundResponse{}).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*getRoundResponse)
		arg.Data = struct {
			Round roundResponse `json:"round"`
		}{
			Round: roundResponse{
				League:      "45",
				Season:      "21",
				Round:       "8",
				IsCompleted: false,
				Games: []gameResponse{
					{
						HomeTeam:  "Rishon",
						GuestTeam: "הפועל קטמון י-ם",
						Stadium:   "Rishon",
						Date:      gt.Format(time.RFC3339),
					},
				},
			},
		}
	})

	r, g, err := GetNextKatamonGame(ctx, &fc)

	if err != nil {
		t.Fatal("Expected nil, got: ", err)
	}

	assert.Equal(t, "45", r.LeagueID, "Expected league 45")
	assert.Equal(t, "21", r.SeasonID, "Expected season 21")
	assert.Equal(t, "8", r.RoundID, "Expected round 8")
	assert.False(t, r.IsCompleted, "Expected false")
	assert.Equal(t, g.HomeTeam, "Rishon", "Expected Rishon to be the home team")
	assert.Equal(t, g.GuestTeam, "הפועל קטמון י-ם", "Expected Katamon to be the guest team")
	assert.Equal(t, g.Stadium, "Rishon", "Expected Rishon to be the stadium")
	assert.Equal(t, g.Date.Format(time.RFC3339), gt.Format(time.RFC3339), "Expected game time to be correct")
}
