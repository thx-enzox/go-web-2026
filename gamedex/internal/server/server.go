package server

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"gamedex/internal/api"
	"gamedex/internal/catalog"
	"gamedex/internal/store"
)

type Server struct {
	Catalog   *catalog.Catalog
	Templates *template.Template
}

type ListPage struct {
	Games  []api.Game
	Genres []string
	Genre  string
	Query  string
	Likes  map[int]int
}

type GamePage struct {
	Game         api.Game
	Tournaments  []api.Tournament
	Likes        int
	Registrations map[int][]store.Registration
}

func New(c *catalog.Catalog) *Server {
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))

	return &Server{
		Catalog:   c,
		Templates: tmpl,
	}
}

func (s *Server) render(w http.ResponseWriter, name string, data any) {
	err := s.Templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println("Erreur template :", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
	}
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/jeux", http.StatusSeeOther)
}

func (s *Server) listGames(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	genre := r.URL.Query().Get("genre")

	games, err := api.SearchGames(q, genre)
	if err != nil {
		log.Println(err)
		http.Error(w, "Le service des jeux ne répond pas.", http.StatusBadGateway)
		return
	}

	genres, err := api.FetchGenres()
	if err != nil {
		log.Println(err)
		http.Error(w, "Impossible de récupérer les genres.", http.StatusBadGateway)
		return
	}

	// Bonus : classement par likes.
	sort.SliceStable(games, func(i, j int) bool {
		return store.Likes(games[i].ID) > store.Likes(games[j].ID)
	})

	s.render(w, "liste.html", ListPage{
		Games:  games,
		Genres: genres,
		Genre:  genre,
		Query:  q,
		Likes:  store.AllLikes(),
	})
}

func (s *Server) gameDetail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	game, ok := s.Catalog.FindGame(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	tournaments := s.Catalog.TournamentsForGame(id)

	registrations := make(map[int][]store.Registration)

	for _, tournament := range tournaments {
		registrations[tournament.ID] = store.Registrations(id, tournament.ID)
	}

	s.render(w, "jeu.html", GamePage{
		Game:          game,
		Tournaments:   tournaments,
		Likes:         store.Likes(id),
		Registrations: registrations,
	})
}

func (s *Server) likeGame(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	_, ok := s.Catalog.FindGame(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	store.AddLike(id)

	// Redirection après le POST : F5 ne refait pas le POST.
	http.Redirect(w, r, "/jeux/"+strconv.Itoa(id), http.StatusSeeOther)
}

func (s *Server) registerTournament(w http.ResponseWriter, r *http.Request) {
	gameID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	tournamentID, err := strconv.Atoi(r.PathValue("tournamentID"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	_, ok := s.Catalog.FindGame(gameID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	tournaments := s.Catalog.TournamentsForGame(gameID)

	tournamentExists := false
	for _, tournament := range tournaments {
		if tournament.ID == tournamentID {
			tournamentExists = true
			break
		}
	}

	if !tournamentExists {
		http.NotFound(w, r)
		return
	}

	pseudo := strings.TrimSpace(r.FormValue("pseudo"))

	if pseudo != "" && len(pseudo) <= 30 {
		store.Register(gameID, tournamentID, pseudo)
	}

	http.Redirect(w, r, "/jeux/"+strconv.Itoa(gameID), http.StatusSeeOther)
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", s.home)
	mux.HandleFunc("GET /jeux", s.listGames)
	mux.HandleFunc("GET /jeux/{id}", s.gameDetail)
	mux.HandleFunc("POST /jeux/{id}/like", s.likeGame)
	mux.HandleFunc("POST /jeux/{id}/tournoi/{tournamentID}/inscription", s.registerTournament)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	return mux
}
