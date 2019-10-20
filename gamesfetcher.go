package main

import (
	"context"
	"time"

	"github.com/shurcooL/graphql"
)

const katamon = "הפועל קטמון י-ם"

// Game contains data about a single game
type Game struct {
	Date      time.Time
	HomeTeam  string
	GuestTeam string
	Stadium   string
}

// Round contains data about a single season round
type Round struct {
	LeagueID    string
	SeasonID    string
	RoundID     string
	IsCompleted bool
	Games       []Game
}

type graphqlGame struct {
	homeTeam  graphql.String
	guestTeam graphql.String
	stadium   graphql.String
	date      graphqlDateTime
}

type graphqlDateTime struct {
	iso graphql.String
}

type graphqlRound struct {
	league      graphql.String
	season      graphql.String
	round       graphql.String
	isCompleted graphql.Boolean
	games       []graphqlGame
}

type roundQuery struct {
	round graphqlRound
}

type querier interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

func getNextRound(ctx context.Context, q querier) (*Round, error) {
	var query roundQuery
	err := q.Query(ctx, &query, nil)
	if err != nil {
		return nil, err
	}

	var games []Game

	for _, g := range query.round.games {
		t, _ := time.Parse(time.RFC3339, string(g.date.iso))
		games = append(games, Game{
			HomeTeam:  string(g.homeTeam),
			GuestTeam: string(g.guestTeam),
			Stadium:   string(g.stadium),
			Date:      t,
		})
	}

	r := Round{
		LeagueID:    string(query.round.league),
		SeasonID:    string(query.round.season),
		RoundID:     string(query.round.round),
		IsCompleted: bool(query.round.isCompleted),
		Games:       games,
	}

	return &r, nil
}

// GetNextKatamonGame Finds the next round in which katamon plays and returns the round information and the game information
func GetNextKatamonGame(ctx context.Context, q querier) (*Round, *Game, error) {
	r, err := getNextRound(ctx, q)
	if err != nil {
		return nil, nil, err
	}

	var g Game

	for _, gi := range r.Games {
		if gi.HomeTeam == katamon || gi.GuestTeam == katamon {
			g = gi
			break
		}
	}

	return r, &g, nil
}
