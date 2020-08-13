package wallet

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zzell/transfer/db/repo"
	"github.com/zzell/transfer/model"
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
		render.JSON(w, http.StatusInternalServerError, model.ErrRsp{
			Error:       "invalid request parameter",
			Description: "failed to parse wallet ID",
		})
	}

	wallet, err := h.repo.GetWallet(walletIDInt)
	if err != nil {
		if err == sql.ErrNoRows {
			render.JSON(w, http.StatusNotFound, model.ErrRsp{
				Error:       "wallet does not exist",
				Description: err.Error(),
			})
			return
		}

		render.JSON(w, http.StatusInternalServerError, model.ErrRsp{
			Error:       "failed to fetch wallet",
			Description: err.Error(),
		})
		return
	}

	render.JSON(w, http.StatusOK, wallet)
}
