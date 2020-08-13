package model

// TransferPayload request model
type TransferPayload struct {
	From   int     `json:"from"`
	To     int     `json:"to"`
	Amount float64 `json:"amount"`
}

// Currency database model
type Currency struct {
	ID     int    `json:"-"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

// Wallet database model
type Wallet struct {
	ID       int      `json:"-"`
	Currency Currency `json:"currency"`
	Score    float64  `json:"score"`
}
