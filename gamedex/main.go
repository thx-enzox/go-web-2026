package main

import (
	"fmt"
	"net/http"
)

func main() {
	
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bienvenue sur GameDex !")
	})

	http.HandleFunc("GET /jeux", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Liste des jeux")
	})

	http.HandleFunc("GET /jeux/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "Jeu numéro %s", id)
	})

	fmt.Println("Serveur sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
