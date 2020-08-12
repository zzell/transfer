package render

import (
	"encoding/json"
	"net/http"
)

const (
	contentType     = "Content-Type"
	jsonContentType = "application/json; charset=utf-8"
)

// Status sends status code
func Status(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

// JSON sends json response with status code
func JSON(w http.ResponseWriter, status int, body interface{}) {
	b, err := json.Marshal(body)
	if err != nil {
		Status(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, jsonContentType)
	Status(w, status)
	_, _ = w.Write(b)
}
