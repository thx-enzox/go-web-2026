
package server

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"gamedex/internal/catalog"
	"gamedex/internal/api"
)

type Server struct {
	Catalog  *catalog.Catalog
	Templates *template.Template
}

type ListPage struct {
	Games []api.Game
	Query string
}

type GamePage struct {
	Game        api.Game
	Tournaments []api.Tournament
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
		http.Error(w, "Erreur d'affichage", http.StatusInternalServerError)
	}
}

func (s *Server) showError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	s.render(w, "erreur.html", message)
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/jeux", http.StatusSeeOther)
}

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	games := s.Catalog.Filter(query)

	s.render(w, "liste.html", ListPage{
		Games: games,
		Query: query,
	})
}

func (s *Server) detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		s.showError(w, http.StatusBadRequest, "Identifiant invalide")
		return
	}

	game, found := s.Catalog.FindGame(id)
	if !found {
		s.showError(w, http.StatusNotFound, "Jeu introuvable")
		return
	}

	tournaments := s.Catalog.TournamentsForGame(game.ID)

	s.render(w, "jeu.html", GamePage{
		Game:        game,
		Tournaments: tournaments,
	})
}

func (s *Server) broken(w http.ResponseWriter, r *http.Request) {
	s.showError(w, http.StatusInternalServerError, "Erreur serveur")
}

func (s *Server) slow(w http.ResponseWriter, r *http.Request) {
	s.showError(w, http.StatusBadGateway, "Service indisponible")
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", s.home)
	mux.HandleFunc("GET /jeux", s.list)
	mux.HandleFunc("GET /jeux/{id}", s.detail)

	mux.HandleFunc("GET /api/broken", s.broken)
	mux.HandleFunc("GET /api/slow", s.slow)

	mux.Handle("GET /static/", http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("web/static")),
	))

	return mux
}
