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
	Handler struct {
		repo      repo.Repository
		config    *cfg.Config
		converter currency.Converter
	}

	Response struct {
		Commission string `json:"commission"`
	}
)

func NewHandler(config *cfg.Config, r repo.Repository, converter currency.Converter) Handler {
	return Handler{
		repo:      r,
		config:    config,
		converter: converter,
	}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var payload = new(model.TransferPayload)

	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		renderErr(w, http.StatusBadRequest, "invalid JSON body", err.Error())
		return
	}

	from, to, err := h.wallets(payload.From, payload.To)
	if err != nil {
		if err == sql.ErrNoRows {
			renderErr(w, http.StatusBadRequest, "wallet does not exist", err.Error())
			return
		}
		renderErr(w, http.StatusInternalServerError, "failed to fetch wallet info", err.Error())
		return
	}

	if from.Score < payload.Amount {
		renderErr(w, http.StatusBadRequest, "failed to process", "too small score in sender's wallet")
		return
	}

	commission := commission(payload.Amount, h.config.CommissionPercent)
	addScore := payload.Amount - commission

	if from.Currency.Symbol != to.Currency.Symbol {
		addScore, err = h.converter.Convert(from.Currency, to.Currency, payload.Amount)
		if err != nil {
			renderErr(w, http.StatusInternalServerError, "failed to convert currency", err.Error())
			return
		}
	}

	err = h.repo.WalletsRepo.Transfer(payload.From, payload.To, payload.Amount, addScore)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			renderErr(w, http.StatusNotFound, "wallet does not exist", err.Error())
			return
		}

		renderErr(w, http.StatusInternalServerError, "failed to process transaction", err.Error())
		return
	}

	render.JSON(w, http.StatusOK, Response{Commission: fmt.Sprintf("%v%s", commission, strings.ToUpper(from.Currency.Symbol))})
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

func renderErr(w http.ResponseWriter, status int, err, desc string) {
	render.JSON(w, status, model.ErrRsp{
		Error:       err,
		Description: desc,
	})
}
