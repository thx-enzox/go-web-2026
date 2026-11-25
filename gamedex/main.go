package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", accueil)
	http.HandleFunc("/jeux", afficherJeux)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func accueil(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func afficherJeux(w http.ResponseWriter, r *http.Request) {
	games, err := FetchGames()
	if err != nil {
		http.Error(w, "Impossible de récupérer les jeux : "+err.Error(), http.StatusInternalServerError)
		return
	}

	tournaments, err := FetchTournaments()
	if err != nil {
		http.Error(w, "Impossible de récupérer les tournois : "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "<h1>GameDex - Liste des jeux</h1>")

	for _, game := range games {
		nombreTournois := 0

		for _, tournament := range tournaments {
			if tournament.GameID == game.ID {
				nombreTournois++
			}
		}

		fmt.Fprintf(w, "<h2>%s (%d) - %d tournoi(s)</h2>",
			game.Name, game.Year, nombreTournois)
	}
}
