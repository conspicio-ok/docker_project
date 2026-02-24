// =============================================
// STORE JV - API Go
// =============================================
// Structure des fichiers :
//   main.go
//   db.go
//   models/game.go
//   models/cart.go
//   handlers/games.go
//   handlers/cart.go
//
// go mod init store-jv
// go get github.com/go-sql-driver/mysql
// go get github.com/gorilla/mux
// =============================================

// ─────────────────────────────────────────────
// main.go
// ─────────────────────────────────────────────
package main

import (
	"log"
	"net/http"

	"store-jv/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Initialise la connexion DB (définie dans db.go)
	InitDB()

	r := mux.NewRouter()

	// Middleware CORS appliqué sur toutes les routes
	// Nécessaire pour que le front (autre origine) puisse appeler l'API
	r.Use(corsMiddleware)

	// ── Routes jeux ──────────────────────────
	r.HandleFunc("/games", handlers.GetGames).Methods("GET")
	r.HandleFunc("/games/{id}", handlers.GetGame).Methods("GET")

	// ── Routes panier ────────────────────────
	// {session_id} : UUID généré côté front, passé dans l'URL
	r.HandleFunc("/cart/{session_id}", handlers.GetCart).Methods("GET")
	r.HandleFunc("/cart/{session_id}/add", handlers.AddToCart).Methods("POST")
	r.HandleFunc("/cart/{session_id}/remove/{game_id}", handlers.RemoveFromCart).Methods("DELETE")

	// Gère les requêtes OPTIONS (preflight CORS)
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("API démarrée sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// corsMiddleware ajoute les headers nécessaires pour autoriser les appels
// depuis le front (qui tourne sur un autre port ou origine)
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // En prod : remplacer * par l'URL du front
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	})
}


// ─────────────────────────────────────────────
// db.go
// ─────────────────────────────────────────────

// package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DB est la variable globale de connexion, accessible depuis les handlers
var DB *sql.DB

func InitDB() {
	var err error
	// Format DSN : user:password@tcp(host:port)/dbname?params
	dsn := "root:password@tcp(127.0.0.1:3306)/store_jv?parseTime=true"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erreur connexion DB :", err)
	}
	// Ping vérifie que la connexion est réellement établie
	if err = DB.Ping(); err != nil {
		log.Fatal("DB injoignable :", err)
	}
	log.Println("Connexion DB OK")
}


// ─────────────────────────────────────────────
// models/game.go
// ─────────────────────────────────────────────

// package models

// Game représente un jeu du catalogue
// Les tags json définissent les noms des champs dans la réponse JSON
type Game struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	Stock       int     `json:"stock"`
}


// ─────────────────────────────────────────────
// models/cart.go
// ─────────────────────────────────────────────

// package models

// CartItem représente une ligne de panier : un jeu + sa quantité
type CartItem struct {
	GameID    int     `json:"game_id"`
	Title     string  `json:"title"`
	Price     float64 `json:"price"`
	ImageURL  string  `json:"image_url"`
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"` // price * quantity, calculé en Go
}

// Cart représente un panier complet
type Cart struct {
	SessionID string     `json:"session_id"`
	Items     []CartItem `json:"items"`
	Total     float64    `json:"total"`
}

// AddToCartRequest est le body JSON attendu pour POST /cart/{session_id}/add
type AddToCartRequest struct {
	GameID   int `json:"game_id"`
	Quantity int `json:"quantity"`
}


// ─────────────────────────────────────────────
// handlers/games.go
// ─────────────────────────────────────────────

// package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"store-jv/models"

	"github.com/gorilla/mux"
)

// GetGames retourne tous les jeux du catalogue
// GET /games
func GetGames(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT id, title, description, price, image_url, stock FROM games")
	if err != nil {
		http.Error(w, "Erreur DB", http.StatusInternalServerError)
		return
	}
	defer rows.Close() // Important : libère la connexion après lecture

	var games []models.Game
	for rows.Next() {
		var g models.Game
		// Scan lit une ligne et remplit la struct
		if err := rows.Scan(&g.ID, &g.Title, &g.Description, &g.Price, &g.ImageURL, &g.Stock); err != nil {
			http.Error(w, "Erreur lecture", http.StatusInternalServerError)
			return
		}
		games = append(games, g)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

// GetGame retourne un jeu par son ID
// GET /games/{id}
func GetGame(w http.ResponseWriter, r *http.Request) {
	// mux.Vars extrait les paramètres de l'URL définis avec {id}
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


// ─────────────────────────────────────────────
// handlers/cart.go
// ─────────────────────────────────────────────

// package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"store-jv/models"

	"github.com/gorilla/mux"
)

// getOrCreateCart récupère l'ID du cart pour ce session_id,
// ou le crée s'il n'existe pas encore (pattern "upsert")
func getOrCreateCart(sessionID string) (int, error) {
	var cartID int
	err := DB.QueryRow("SELECT id FROM cart WHERE session_id = ?", sessionID).Scan(&cartID)
	if err != nil {
		// Cart inexistant : on le crée
		result, err := DB.Exec("INSERT INTO cart (session_id) VALUES (?)", sessionID)
		if err != nil {
			return 0, err
		}
		id, _ := result.LastInsertId()
		return int(id), nil
	}
	return cartID, nil
}

// GetCart retourne le contenu du panier
// GET /cart/{session_id}
func GetCart(w http.ResponseWriter, r *http.Request) {
	sessionID := mux.Vars(r)["session_id"]
	cartID, err := getOrCreateCart(sessionID)
	if err != nil {
		http.Error(w, "Erreur cart", http.StatusInternalServerError)
		return
	}

	// JOIN pour récupérer les infos du jeu en même temps
	rows, err := DB.Query(`
		SELECT g.id, g.title, g.price, g.image_url, ci.quantity
		FROM cart_items ci
		JOIN games g ON g.id = ci.game_id
		WHERE ci.cart_id = ?`, cartID)
	if err != nil {
		http.Error(w, "Erreur DB", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	cart := models.Cart{SessionID: sessionID, Items: []models.CartItem{}}
	for rows.Next() {
		var item models.CartItem
		rows.Scan(&item.GameID, &item.Title, &item.Price, &item.ImageURL, &item.Quantity)
		item.Subtotal = item.Price * float64(item.Quantity)
		cart.Total += item.Subtotal
		cart.Items = append(cart.Items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// AddToCart ajoute un jeu au panier (ou incrémente la quantité)
// POST /cart/{session_id}/add
// Body JSON : { "game_id": 1, "quantity": 1 }
func AddToCart(w http.ResponseWriter, r *http.Request) {
	sessionID := mux.Vars(r)["session_id"]

	var req models.AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.GameID == 0 {
		http.Error(w, "Body invalide", http.StatusBadRequest)
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	cartID, err := getOrCreateCart(sessionID)
	if err != nil {
		http.Error(w, "Erreur cart", http.StatusInternalServerError)
		return
	}

	// INSERT ... ON DUPLICATE KEY UPDATE : si le jeu est déjà dans le panier,
	// on additionne la quantité au lieu d'insérer une nouvelle ligne
	_, err = DB.Exec(`
		INSERT INTO cart_items (cart_id, game_id, quantity)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE quantity = quantity + VALUES(quantity)`,
		cartID, req.GameID, req.Quantity)
	if err != nil {
		http.Error(w, "Erreur ajout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// RemoveFromCart supprime un jeu du panier
// DELETE /cart/{session_id}/remove/{game_id}
func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["session_id"]
	gameID, err := strconv.Atoi(vars["game_id"])
	if err != nil {
		http.Error(w, "game_id invalide", http.StatusBadRequest)
		return
	}

	var cartID int
	err = DB.QueryRow("SELECT id FROM cart WHERE session_id = ?", sessionID).Scan(&cartID)
	if err != nil {
		http.Error(w, "Cart introuvable", http.StatusNotFound)
		return
	}

	DB.Exec("DELETE FROM cart_items WHERE cart_id = ? AND game_id = ?", cartID, gameID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "removed"})
}