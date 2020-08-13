package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zzell/transfer/cfg"
	"github.com/zzell/transfer/currency"
	"github.com/zzell/transfer/db/repo"
	"github.com/zzell/transfer/web/transfer"
	"github.com/zzell/transfer/web/wallet"
)

// NewRouter constructor
func NewRouter(repository repo.Repository, config *cfg.Config, converter currency.Converter) *mux.Router {
	r := mux.NewRouter()

	transferHandler := transfer.NewHandler(config, repository, converter)
	walletHandler := wallet.NewHandler(repository)

	r.HandleFunc("/transfer", transferHandler.Handle).Methods(http.MethodPost)
	r.HandleFunc("/wallet/{id}", walletHandler.GetWallet).Methods(http.MethodGet)

	return r
}
