
package main

import (
	"log"
	"net/http"

	"gamedex/internal/api"
	"gamedex/internal/catalog"
	"gamedex/internal/server"
)

func main() {
	// Chargement des jeux une seule fois au démarrage.
	games, err := api.FetchGames()
	if err != nil {
		log.Fatal("Impossible de charger les jeux :", err)
	}

	// Chargement des tournois une seule fois au démarrage.
	tournaments, err := api.FetchTournaments()
	if err != nil {
		log.Fatal("Impossible de charger les tournois :", err)
	}

	// Création du catalogue en mémoire.
	c := catalog.New(games, tournaments)

	// Création du serveur web.
	s := server.New(c)

	// Lancement du serveur.
	log.Println("GameDex lancé sur http://localhost:8080")

	err = http.ListenAndServe(":8080", s.Routes())
	if err != nil {
		log.Fatal(err)
	}
}
