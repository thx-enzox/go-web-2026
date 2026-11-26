
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Game struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Year int    `json:"year"`
}

type GamePage struct {
	Game        Game
	Tournaments []Tournament
}

type ListPage struct {
	Games []Game
	Query string
}

var games []Game
var tournaments []Tournament

var templates = template.Must(template.ParseGlob("templates/*.html"))

func render(w http.ResponseWriter, name string, data any) {
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println("Erreur template :", err)
		http.Error(w, "Erreur d'affichage", http.StatusInternalServerError)
	}
}

func showError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	render(w, "erreur.html", message)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		showError(w, http.StatusNotFound, "Page introuvable")
		return
	}
	http.Redirect(w, r, "/jeux", http.StatusSeeOther)
}

func gamesHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/jeux" {
		showError(w, http.StatusNotFound, "Page introuvable")
		return
	}

	q := r.URL.Query().Get("q")
	filtered := []Game{}

	for _, game := range games {
		if strings.Contains(strings.ToLower(game.Name), strings.ToLower(q)) {
			filtered = append(filtered, game)
		}
	}

	render(w, "liste.html", ListPage{
		Games: filtered,
		Query: q,
	})
}

func gameDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		showError(w, http.StatusBadRequest, "Identifiant invalide")
		return
	}

	var selectedGame Game
	found := false

	for _, game := range games {
		if game.ID == id {
			selectedGame = game
			found = true
			break
		}
	}

	if !found {
		showError(w, http.StatusNotFound, "Jeu introuvable")
		return
	}

	var gameTournaments []Tournament

	for _, tournament := range tournaments {
		if tournament.GameID == selectedGame.ID {
			gameTournaments = append(gameTournaments, tournament)
		}
	}

	render(w, "jeu.html", GamePage{
		Game:        selectedGame,
		Tournaments: gameTournaments,
	})
}

func brokenHandler(w http.ResponseWriter, r *http.Request) {
	showError(w, http.StatusInternalServerError, "Erreur serveur")
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	showError(w, http.StatusBadGateway, "Service indisponible")
}

func main() {
	var err error

	games, err = FetchGames()
	if err != nil {
		log.Fatal("Impossible de charger les jeux :", err)
	}

	tournaments, err = FetchTournaments()
	if err != nil {
		log.Fatal("Impossible de charger les tournois :", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", homeHandler)
	mux.HandleFunc("GET /jeux", gamesHandler)
	mux.HandleFunc("GET /jeux/{id}", gameDetailHandler)

	mux.HandleFunc("GET /api/broken", brokenHandler)
	mux.HandleFunc("GET /api/slow", slowHandler)

	mux.Handle("GET /static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))

	fmt.Println("Serveur lancé sur http://localhost:8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
