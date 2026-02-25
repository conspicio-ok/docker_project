package models

type CartItem struct {
	GameID    int     `json:"game_id"`
	Title     string  `json:"title"`
	Price     float64 `json:"price"`
	ImageURL  string  `json:"image_url"`
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
}

type Cart struct {
	SessionID string     `json:"session_id"`
	Items     []CartItem `json:"items"`
	Total     float64    `json:"total"`
}

type AddToCartRequest struct {
	GameID   int `json:"game_id"`
	Quantity int `json:"quantity"`
}