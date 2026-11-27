package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"
)

type Game struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Year  int    `json:"year"`
	Genre string `json:"genre"`
}

type Tournament struct {
	ID     int    `json:"id"`
	GameID int    `json:"gameId"`
	Name   string `json:"name"`
	City   string `json:"city"`
	Date   string `json:"date"`
}

const baseURL = "http://localhost:8081"

var client = &http.Client{
	Timeout: 5 * time.Second,
}

func SearchGames(q, genre string) ([]Game, error) {
	params := url.Values{}
	params.Set("q", q)
	params.Set("genre", genre)

	resp, err := client.Get(baseURL + "/api/games?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API games: statut %d", resp.StatusCode)
	}

	var games []Game

	if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
		return nil, err
	}

	return games, nil
}

func FetchGames() ([]Game, error) {
	return SearchGames("", "")
}

func FetchTournaments() ([]Tournament, error) {
	resp, err := client.Get(baseURL + "/api/tournaments")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API tournaments: statut %d", resp.StatusCode)
	}

	var tournaments []Tournament

	if err := json.NewDecoder(resp.Body).Decode(&tournaments); err != nil {
		return nil, err
	}

	return tournaments, nil
}

func FetchGenres() ([]string, error) {
	games, err := FetchGames()
	if err != nil {
		return nil, err
	}

	unique := make(map[string]bool)

	for _, game := range games {
		if game.Genre != "" {
			unique[game.Genre] = true
		}
	}

	genres := make([]string, 0, len(unique))

	for genre := range unique {
		genres = append(genres, genre)
	}

	sort.Strings(genres)

	return genres, nil
}
