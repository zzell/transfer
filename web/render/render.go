package render

import (
	"encoding/json"
	"net/http"
)

const (
	contentType     = "Content-Type"
	jsonContentType = "application/json; charset=utf-8"
)

// ErrRsp response model
type ErrRsp struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}

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

// Error sends error response
func Error(w http.ResponseWriter, status int, err, desc string) {
	JSON(w, status, ErrRsp{
		Error:       err,
		Description: desc,
	})
}
