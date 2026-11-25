package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Tournament struct {
	ID     int    `json:"id"`
	GameID int    `json:"game_id"`
	Name   string `json:"name"`
}

var client = http.Client{
	Timeout: 5 * time.Second,
}

func FetchGames() ([]Game, error) {
	resp, err := client.Get("http://localhost:8081/api/games")
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
	resp, err := client.Get("http://localhost:8081/api/tournaments")
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

func afficherJeux() error {
	games, err := FetchGames()
	if err != nil {
		return err
	}

	tournaments, err := FetchTournaments()
	if err != nil {
		return err
	}

	fmt.Println("\n===== GAMEDEX =====")

	for _, game := range games {
		nombreTournois := 0

		for _, tournament := range tournaments {
			if tournament.GameID == game.ID {
				nombreTournois++
			}
		}

		fmt.Printf("%s (%d) - %d tournoi(s)\n",
			game.Name, game.Year, nombreTournois)
	}

	return nil
}
