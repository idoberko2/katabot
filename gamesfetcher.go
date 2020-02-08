package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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

const getRoundQuery = `
query getRound($league: String!,$season: String!,$round: String!) { 
	round(league: $league, season: $season, round: $round) { 
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

type getRoundResponse struct {
	Data struct {
		Round roundResponse `json:"round"`
	} `json:"data"`
}

func extractRoundFromResponse(r *roundResponse) *Round {
	var games []Game

	for _, g := range r.Games {
		t, err := time.Parse(time.RFC3339, string(g.Date))
		if err != nil {
			log.Print("Error when parsing time: ", err)
		}

		if os.Getenv("MATCHES_TIMEZONE") != "" {
			loc, err := time.LoadLocation(os.Getenv("MATCHES_TIMEZONE"))
			if err != nil {
				log.Print("Error loading timezone: ", err)
			}
			t = t.In(loc)
		}

		games = append(games, Game{
			HomeTeam:  g.HomeTeam,
			GuestTeam: g.GuestTeam,
			Stadium:   g.Stadium,
			Date:      t,
		})
	}

	rr := Round{
		LeagueID:    r.League,
		SeasonID:    r.Season,
		RoundID:     r.Round,
		IsCompleted: r.IsCompleted,
		Games:       games,
	}

	return &rr
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

	return extractRoundFromResponse(&j.Data.NextRound), nil
}

func getSpecificRound(ctx context.Context, c apiRequestMaker, l, s, r string) (*Round, error) {
	var j getRoundResponse
	err := c.Request(map[string]interface{}{
		"operationName": "getRound",
		"query":         getRoundQuery,
		"variables": map[string]string{
			"league": l,
			"season": s,
			"round":  r,
		},
	}, &j)
	if err != nil {
		return nil, err
	}

	return extractRoundFromResponse(&j.Data.Round), nil
}

func findKatamonGame(r *Round) *Game {
	for _, g := range r.Games {
		if g.HomeTeam == katamon || g.GuestTeam == katamon {
			return &g
		}
	}

	return nil
}

// GamesFetcher is an interface for fetching Katamon games
type GamesFetcher interface {
	GetNextKatamonGame(ctx context.Context) (*Round, *Game, error)
}

// GamesFetcherInst is a wrapper for fetching Katamon games
type GamesFetcherInst struct {
	Client apiRequestMaker
}

// GetNextKatamonGame Finds the next round in which katamon plays and returns the round information and the game information
func (gf *GamesFetcherInst) GetNextKatamonGame(ctx context.Context) (*Round, *Game, error) {
	r, err := getNextRound(ctx, gf.Client)
	if err != nil {
		return nil, nil, err
	}

	g := findKatamonGame(r)

	if g.Date.Before(time.Now()) {
		ri, err := strconv.Atoi(r.RoundID)
		if err != nil {
			return nil, nil, err
		}
		r, err = getSpecificRound(ctx, gf.Client, r.LeagueID, r.SeasonID, strconv.Itoa(ri+1))
		if err != nil {
			return nil, nil, err
		}

		g = findKatamonGame(r)
	}

	return r, g, nil
}
