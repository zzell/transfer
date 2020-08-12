package transfer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zzell/transfer/db/repo"
	"github.com/zzell/transfer/model"
	"github.com/zzell/transfer/web/render"
)

type (
	Handler struct {
		walletsRepo repo.WalletsRepo
	}

	payload struct {
		From   int     `json:"from"`
		To     int     `json:"to"`
		Amount float64 `json:"amount"`
	}
)

func NewHandler(r repo.WalletsRepo) Handler {
	return Handler{walletsRepo: r}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var body = new(payload)

	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		renderErr(w, http.StatusBadRequest, "invalid JSON body", err.Error())
		return
	}

	senderScore, err := h.walletsRepo.GetScore(body.From)
	if err != nil {
		if err == sql.ErrNoRows {
			renderErr(w, http.StatusBadRequest, "wallet does not exist", err.Error())
			return
		}

		renderErr(w, http.StatusInternalServerError, "failed to connect to external resource", err.Error())
		return
	}

	if senderScore < body.Amount {
		renderErr(w, http.StatusBadRequest, "failed to process", "not enough score in sender's wallet")
		return
	}

	err = h.walletsRepo.Transfer(body.From, body.To, body.Amount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			renderErr(w, http.StatusNotFound, "wallet does not exist", err.Error())
			return
		}

		renderErr(w, http.StatusInternalServerError, "failed to process transaction", err.Error())
		return
	}

	render.Status(w, http.StatusOK)
}

func renderErr(w http.ResponseWriter, status int, err, desc string) {
	render.JSON(w, status, model.Error{
		Error:       err,
		Description: desc,
	})
}
