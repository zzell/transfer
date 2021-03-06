package transfer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zzell/transfer/cfg"
	"github.com/zzell/transfer/currency"
	"github.com/zzell/transfer/db/repo"
	"github.com/zzell/transfer/model"
	"github.com/zzell/transfer/web/render"
)

type (
	// Handler dependency bridge
	Handler struct {
		repo      repo.Repository
		config    *cfg.Config
		converter currency.Converter
	}

	response struct {
		Commission string `json:"commission"`
	}
)

// NewHandler constructor
func NewHandler(config *cfg.Config, r repo.Repository, converter currency.Converter) Handler {
	return Handler{
		repo:      r,
		config:    config,
		converter: converter,
	}
}

// Transfer handles transfer request
func (h Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	var payload = new(model.TransferPayload)

	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "invalid JSON body", err.Error())
		return
	}

	from, to, err := h.wallets(payload.From, payload.To)
	if errors.Is(err, sql.ErrNoRows) {
		render.Error(w, http.StatusBadRequest, "wallet does not exist", err.Error())
		return
	}

	if err != nil {
		render.Error(w, http.StatusInternalServerError, "failed to fetch wallet info", err.Error())
		return
	}

	if from.Score < payload.Amount {
		render.Error(w, http.StatusBadRequest, "failed to process", "too small score in sender's wallet")
		return
	}

	commission := commission(payload.Amount, h.config.CommissionPercent)
	addScore := payload.Amount - commission

	if from.Currency.Symbol != to.Currency.Symbol {
		addScore, err = h.converter.Convert(from.Currency, to.Currency, payload.Amount)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, "failed to convert currency", err.Error())
			return
		}
	}

	err = h.repo.WalletsRepo.Transfer(payload.From, payload.To, payload.Amount, addScore)
	if errors.Is(err, sql.ErrNoRows) {
		render.Error(w, http.StatusNotFound, "wallet does not exist", err.Error())
		return
	}

	if err != nil {
		render.Error(w, http.StatusInternalServerError, "failed to process transaction", err.Error())
		return
	}

	render.JSON(w, http.StatusOK, response{Commission: fmt.Sprintf("%v%s", commission, strings.ToUpper(from.Currency.Symbol))})
}

// fetches two wallets concurrently
func (h Handler) wallets(fromID, toID int) (from, to *model.Wallet, err error) {
	var (
		goroutines = 2
		errch      = make(chan error, goroutines)
	)

	go func() {
		var e error
		from, e = h.repo.WalletsRepo.GetWallet(fromID)
		errch <- e
	}()

	go func() {
		var e error
		to, e = h.repo.WalletsRepo.GetWallet(toID)
		errch <- e
	}()

	for i := 0; i < goroutines; i++ {
		err = <-errch
		if err != nil {
			return
		}
	}

	return
}

// get commission from gross value
func commission(gross, percentage float64) float64 {
	return (percentage / 100) * gross
}
