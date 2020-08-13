package currency

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/zzell/transfer/model"
)

//go:generate mockgen -package mock -destination=mock/converter_mock.go github.com/zzell/transfer/currency Converter

// Converter converts value of one currency to another
type Converter interface {
	Convert(from, to model.Currency, amount float64) (float64, error)
}

const coingeckoPriceAPI = "https://api.coingecko.com/api/v3/simple/price"

var errInvalidRsp = "invalid response from external service"

// CoingeckoConverter uses third-party API to fetch relevant cryptocurrency relations
type (
	CoingeckoConverter struct {
		// keep http client to reuse TCP connection
		Client *http.Client
	}

	coingeckoRsp map[string]map[string]float64
)

// NewConverter constructor
// uses robustness principle
func NewConverter() *CoingeckoConverter {
	return &CoingeckoConverter{
		// default client for simplicity, in real life custom one should be used
		Client: http.DefaultClient,
	}
}

// Convert converts currency using relation from API
func (c *CoingeckoConverter) Convert(from, to model.Currency, amount float64) (float64, error) {
	req, err := http.NewRequest(http.MethodGet, coingeckoPriceAPI, nil)
	if err != nil {
		return 0, err
	}

	// for some reason their API consumes NAME as base and SYMBOL as relative
	values := url.Values{
		"ids":           []string{from.Name},
		"vs_currencies": []string{to.Symbol},
	}

	req.URL.RawQuery = values.Encode()

	rsp, err := c.Client.Do(req)
	if err != nil {
		return 0, err
	}

	defer func() { _ = rsp.Body.Close() }()

	if rsp.StatusCode != http.StatusOK {
		return 0, errors.New(errInvalidRsp)
	}

	var body = make(coingeckoRsp)

	err = json.NewDecoder(rsp.Body).Decode(&body)
	if err != nil {
		return 0, err
	}

	relations, ok := body[from.Name]
	if !ok {
		return 0, errors.New(errInvalidRsp)
	}

	value, ok := relations[to.Symbol]
	if !ok {
		return 0, errors.New(errInvalidRsp)
	}

	return amount * value, nil
}
