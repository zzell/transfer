package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zzell/transfer/cfg"
	"github.com/zzell/transfer/db/repo"
	"github.com/zzell/transfer/web/transfer"
	"github.com/zzell/transfer/web/wallet"
)

func NewRouter(repo repo.WalletsRepo, config *cfg.Config) *mux.Router {
	r := mux.NewRouter()

	transferHandler := transfer.NewHandler(repo, config)
	walletHandler := wallet.NewHandler(repo)

	r.HandleFunc("/transfer", transferHandler.Handle).Methods(http.MethodPost)
	r.HandleFunc("/wallet/{id}", walletHandler.GetScore).Methods(http.MethodGet)

	return r
}
