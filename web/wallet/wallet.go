package wallet

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zzell/transfer/db/repo"
	"github.com/zzell/transfer/web/render"
)

const idPathParam = "id"

// Handler dependency bridge
type Handler struct {
	repo repo.Repository
}

// NewHandler constructor
func NewHandler(r repo.Repository) Handler {
	return Handler{repo: r}
}

// GetWallet renders wallet info
func (h Handler) GetWallet(w http.ResponseWriter, r *http.Request) {
	walletID := mux.Vars(r)[idPathParam]

	walletIDInt, err := strconv.Atoi(walletID)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, "invalid request parameter", "failed to parse wallet ID")
		return
	}

	wallet, err := h.repo.GetWallet(walletIDInt)
	if errors.Is(err, sql.ErrNoRows) {
		render.Error(w, http.StatusNotFound, "wallet does not exist", err.Error())
		return
	}

	if err != nil {
		render.Error(w, http.StatusInternalServerError, "failed to fetch wallet", err.Error())
		return
	}

	render.JSON(w, http.StatusOK, wallet)
}
