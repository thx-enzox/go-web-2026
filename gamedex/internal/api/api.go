
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Game struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Year int    `json:"year"`
}

type Tournament struct {
	ID     int    `json:"id"`
	GameID int    `json:"game_id"`
	Name   string `json:"name"`
	City   string `json:"city"`
	Date   string `json:"date"`
}

const baseURL = "http://localhost:8081"

var client = http.Client{
	Timeout: 5 * time.Second,
}

func FetchGames() ([]Game, error) {
	resp, err := client.Get(baseURL + "/api/games")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API games : statut %d", resp.StatusCode)
	}

	var games []Game
	err = json.NewDecoder(resp.Body).Decode(&games)
	if err != nil {
		return nil, err
	}

	return games, nil
}

func FetchTournaments() ([]Tournament, error) {
	resp, err := client.Get(baseURL + "/api/tournaments")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API tournaments : statut %d", resp.StatusCode)
	}

	var tournaments []Tournament
	err = json.NewDecoder(resp.Body).Decode(&tournaments)
	if err != nil {
		return nil, err
	}

	return tournaments, nil
}
