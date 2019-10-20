package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const katamon = "הפועל קטמון י-ם"
const nextRoundQuery = `
query getNextRound($timestamp: Int!) { 
	nextRound(timestamp: $timestamp) { 
		league 
		season 
		round 
		isCompleted 
		games { 
			date 
			homeTeam 
			guestTeam 
			stadium 
		} 
	} 
}`

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

type apiRequestMaker interface {
	Request(map[string]interface{}, interface{}) error
}

// APIClient is the client for interacting with the GraphQL server
type APIClient struct {
	URL  string
	Mime string
}

// Request makes a GraphQL using the supplied query q and populates the response in j
func (c APIClient) Request(q map[string]interface{}, j interface{}) error {
	p, err := json.Marshal(q)
	if err != nil {
		return err
	}
	resp, err := http.Post(c.URL, c.Mime, bytes.NewBuffer(p))
	if err != nil {
		return fmt.Errorf("API request failed with error: %s", err)
	}

	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Invalid API response, failed parsing JSON with error: %s", err)
	}

	err = json.Unmarshal(rb, j)
	if err != nil {
		return fmt.Errorf("Invalid API response, failed parsing JSON with error: %s", err)
	}

	return nil
}

type nextRoundParam struct {
	timestamp int64
}

type gameResponse struct {
	Date      string `json:"date"`
	HomeTeam  string `json:"homeTeam"`
	GuestTeam string `json:"guestTeam"`
	Stadium   string `json:"stadium"`
}

type roundResponse struct {
	League      string         `json:"league"`
	Season      string         `json:"season"`
	Round       string         `json:"round"`
	IsCompleted bool           `json:"isCompleted"`
	Games       []gameResponse `json:"games"`
}

type nextRoundResponse struct {
	Data struct {
		NextRound roundResponse `json:"nextRound"`
	} `json:"data"`
}

func getNextRound(ctx context.Context, c apiRequestMaker) (*Round, error) {
	var j nextRoundResponse
	err := c.Request(map[string]interface{}{
		"operationName": "getNextRound",
		"query":         nextRoundQuery,
		"variables": map[string]int64{
			"timestamp": time.Now().Unix(),
		},
	}, &j)
	if err != nil {
		return nil, err
	}

	var games []Game

	for _, g := range j.Data.NextRound.Games {
		t, _ := time.Parse(time.RFC3339, string(g.Date))
		games = append(games, Game{
			HomeTeam:  g.HomeTeam,
			GuestTeam: g.GuestTeam,
			Stadium:   g.Stadium,
			Date:      t,
		})
	}

	r := Round{
		LeagueID:    j.Data.NextRound.League,
		SeasonID:    j.Data.NextRound.Season,
		RoundID:     j.Data.NextRound.Round,
		IsCompleted: j.Data.NextRound.IsCompleted,
		Games:       games,
	}

	return &r, nil
}

// GetNextKatamonGame Finds the next round in which katamon plays and returns the round information and the game information
func GetNextKatamonGame(ctx context.Context, c apiRequestMaker) (*Round, *Game, error) {
	r, err := getNextRound(ctx, c)
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
