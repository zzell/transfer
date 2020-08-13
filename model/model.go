package model

// TransferPayload request model
type TransferPayload struct {
	From   int     `json:"from"`
	To     int     `json:"to"`
	Amount float64 `json:"amount"`
}

// ErrRsp response model
type ErrRsp struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}

// Currency database model
type Currency struct {
	ID     int
	Name   string
	Symbol string
}

// Wallet database model
type Wallet struct {
	ID       int
	Currency Currency
	Score    float64
}
