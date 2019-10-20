package main

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/shurcooL/graphql"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)

type fakeQuerier struct {
	mock.Mock
}

func (fq *fakeQuerier) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	args := fq.Called(ctx, q, variables)

	return args.Error(0)
}

func TestGetNextKatamonGame_error(t *testing.T) {
	fq := fakeQuerier{}
	ctx := context.Background()
	fq.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Some error occurred"))

	_, _, err := GetNextKatamonGame(ctx, &fq)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGetNextKatamonGame_simple(t *testing.T) {
	var query roundQuery
	fq := fakeQuerier{}
	ctx := context.Background()
	gt := time.Now().Add(time.Hour * 12)
	fq.On("Query", ctx, &query, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		fmt.Println(args)
		arg := args.Get(1).(*roundQuery)
		arg.round = graphqlRound{
			league:      "45",
			season:      "21",
			round:       "7",
			isCompleted: false,
			games: []graphqlGame{
				{
					homeTeam:  "הפועל קטמון י-ם",
					guestTeam: "Beitar",
					stadium:   "Teddy",
					date: graphqlDateTime{
						iso: graphql.String(gt.Format(time.RFC3339)),
					},
				},
			},
		}
	})

	r, g, err := GetNextKatamonGame(ctx, &fq)

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
