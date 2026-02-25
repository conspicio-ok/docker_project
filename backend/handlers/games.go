package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"store-jv/models"
	"github.com/gorilla/mux"
)

func GetGames(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT id, title, description, price, image_url, stock FROM games")
	if err != nil {
		http.Error(w, "Erreur DB", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var games []models.Game
	for rows.Next() {
		var g models.Game
		if err := rows.Scan(&g.ID, &g.Title, &g.Description, &g.Price, &g.ImageURL, &g.Stock); err != nil {
			http.Error(w, "Erreur lecture", http.StatusInternalServerError)
			return
		}
		games = append(games, g)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

func GetGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}
	var g models.Game
	err = DB.QueryRow(
		"SELECT id, title, description, price, image_url, stock FROM games WHERE id = ?", id,
	).Scan(&g.ID, &g.Title, &g.Description, &g.Price, &g.ImageURL, &g.Stock)

	if err == sql.ErrNoRows {
		http.Error(w, "Jeu introuvable", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erreur DB", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}