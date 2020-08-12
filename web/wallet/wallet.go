package wallet

import (
	"net/http"

	"github.com/zzell/transfer/db/repo"
)

type (
	Handler struct {
		walletsRepo repo.WalletsRepo
	}
)

func NewHandler(r repo.WalletsRepo) Handler {
	return Handler{walletsRepo: r}
}


func (h Handler) GetScore(w http.ResponseWriter, r *http.Request) {
	// h.GetScore()
}
