package model

// Error response model
type Error struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}
