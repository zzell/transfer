package transfer

import (
	"encoding/json"
	"net/http"

	"github.com/zzell/transfer/model"
	"github.com/zzell/transfer/web/render"
)

type (
	Handler struct {
	}

	payload struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount int    `json:"amount"`
	}
)

func NewHandler() Handler {
	return Handler{}
}

func (Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var body = new(payload)

	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, model.Error{
			Error:       "invalid JSON body",
			Description: err.Error(),
		})
	}
}
