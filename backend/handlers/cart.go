package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"store-jv/models"
	"github.com/gorilla/mux"
)

func getOrCreateCart(sessionID string) (int, error) {
	var cartID int
	err := DB.QueryRow("SELECT id FROM cart WHERE session_id = ?", sessionID).Scan(&cartID)
	if err != nil {
		result, err := DB.Exec("INSERT INTO cart (session_id) VALUES (?)", sessionID)
		if err != nil {
			return 0, err
		}
		id, _ := result.LastInsertId()
		return int(id), nil
	}
	return cartID, nil
}

func GetCart(w http.ResponseWriter, r *http.Request) {
	sessionID := mux.Vars(r)["session_id"]
	cartID, err := getOrCreateCart(sessionID)
	if err != nil {
		http.Error(w, "Erreur cart", http.StatusInternalServerError)
		return
	}
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