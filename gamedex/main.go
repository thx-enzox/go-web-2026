
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Game struct {
	ID   int
	Name string
	Year int
}

var games = []Game{
	{1, "Celeste", 2018},
	{2, "Hades", 2020},
	{3, "Rocket League", 2015},
	{4, "Street Fighter 6", 2023},
	{5, "Stardew Valley", 2016},
	{6, "League of Legends", 2009},
	{7, "Valorant", 2020},
	{8, "Dead Cells", 2018},
	{9, "Trackmania", 2020},
	{10, "Mario Kart 8 Deluxe", 2017},
	{11, "Hollow Knight", 2017},
	{12, "Minecraft", 2011},
}

type PageData struct {
	Games []Game
}

var pageTemplate = template.Must(template.New("games").Parse(`
<!DOCTYPE html>
<html lang="fr">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>GameDex - Jeux</title>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
	<header class="navbar">
		<div class="logo">Game<span>Dex</span></div>
		<nav>
			<a href="/">Accueil</a>
			<a href="/jeux">Jeux</a>
		</nav>
	</header>

	<section class="hero">
		<p class="tag">GAMING DATABASE</p>
		<h1>Ma collection <span>de jeux</span></h1>
		<p>Retrouve tous les jeux disponibles dans ta GameDex.</p>
	</section>

	<main class="container">
		<div class="section-title">
			<div>
				<p class="tag">COLLECTION</p>
				<h2>Liste des jeux</h2>
			</div>
			<div class="count">{{len .Games}} jeux</div>
		</div>

		<div class="games-grid">
			{{range .Games}}
			<div class="game-card">
				<div class="game-id">GAME #{{.ID}}</div>
				<h3>{{.Name}}</h3>
				<p class="year">Sortie : {{.Year}}</p>
				<a class="details" href="/jeux/{{.ID}}">Voir les détails →</a>
			</div>
			{{end}}
		</div>
	</main>

	<footer>
		<p>GameDex © 2026</p>
	</footer>
</body>
</html>
`))

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintln(w, "Bienvenue sur GameDex")
}

func gamesHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/jeux" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := pageTemplate.Execute(w, PageData{Games: games})
	if err != nil {
		http.Error(w, "Erreur d'affichage", http.StatusInternalServerError)
	}
}

func gameDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	if id < 1 || id > len(games) {
		http.Error(w, "Jeu introuvable", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Jeu : %s (%d)", games[id-1].Name, games[id-1].Year)
}

func main() {
	http.HandleFunc("GET /", homeHandler)
	http.HandleFunc("GET /jeux", gamesHandler)
	http.HandleFunc("GET /jeux/{id}", gameDetailHandler)

	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Serveur lancé sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
