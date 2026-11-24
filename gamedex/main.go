package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bienvenue sur GameDex")
	})


	http.HandleFunc("GET /jeux", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Voici la liste des jeux")
	})


	http.HandleFunc("GET /jeux/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		numero, err := strconv.Atoi(id)

		if err != nil {
			http.Error(w, "Identifiant invalide", http.StatusBadRequest)
			return
		}

		if numero < 1 || numero > 3 {
			http.Error(w, "Jeu introuvable", http.StatusNotFound)
			return
		}

		fmt.Fprintf(w, "Vous consultez le jeu numéro %d", numero)
	})

	fmt.Println("Serveur lancé sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
