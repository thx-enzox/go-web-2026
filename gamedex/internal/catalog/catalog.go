
package catalog

import (
	"strings"

	"gamedex/internal/api"
)

type Catalog struct {
	Games       []api.Game
	Tournaments []api.Tournament
}

func New(games []api.Game, tournaments []api.Tournament) *Catalog {
	return &Catalog{
		Games:       games,
		Tournaments: tournaments,
	}
}

func (c *Catalog) Filter(query string) []api.Game {
	result := []api.Game{}

	for _, game := range c.Games {
		if strings.Contains(strings.ToLower(game.Name), strings.ToLower(query)) {
			result = append(result, game)
		}
	}

	return result
}

func (c *Catalog) FindGame(id int) (api.Game, bool) {
	for _, game := range c.Games {
		if game.ID == id {
			return game, true
		}
	}

	return api.Game{}, false
}

func (c *Catalog) TournamentsForGame(gameID int) []api.Tournament {
	result := []api.Tournament{}

	for _, tournament := range c.Tournaments {
		if tournament.GameID == gameID {
			result = append(result, tournament)
		}
	}

	return result
}
