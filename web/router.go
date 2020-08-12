package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zzell/transfer/db/repo"
	"github.com/zzell/transfer/web/transfer"
	"github.com/zzell/transfer/web/wallet"
)

func NewRouter(repo repo.WalletsRepo) *mux.Router {
	r := mux.NewRouter()

	transferHandler := transfer.NewHandler(repo)
	walletHandler := wallet.NewHandler(repo)

	r.HandleFunc("/transfer", transferHandler.Handle).Methods(http.MethodPost)
	r.HandleFunc("/wallet/{id}", walletHandler.GetScore).Methods(http.MethodGet)

	return r
}
