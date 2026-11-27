package store

import "sync"

type Registration struct {
	GameID      int
	TournamentID int
	Pseudo       string
}

var (
	registrationMu sync.Mutex
	registrations  = []Registration{}
)

func Register(gameID, tournamentID int, pseudo string) bool {
	registrationMu.Lock()
	defer registrationMu.Unlock()

	for _, r := range registrations {
		if r.GameID == gameID &&
			r.TournamentID == tournamentID &&
			r.Pseudo == pseudo {
			return false
		}
	}

	registrations = append(registrations, Registration{
		GameID:      gameID,
		TournamentID: tournamentID,
		Pseudo:       pseudo,
	})

	return true
}

func Registrations(gameID, tournamentID int) []Registration {
	registrationMu.Lock()
	defer registrationMu.Unlock()

	var result []Registration

	for _, r := range registrations {
		if r.GameID == gameID && r.TournamentID == tournamentID {
			result = append(result, r)
		}
	}

	return result
}
