package repo

import (
	"database/sql"

	"github.com/zzell/transfer/model"
)

const (
	selectCurrencyListSQL = `SELECT id, symbol, name FROM currencies`
)

type (
	CurrenciesRepo interface {
		GetList() ([]model.Currency, error)
	}

	currencies struct {
		conn *sql.DB
	}
)

func NewCurrenciesRepo(conn *sql.DB) CurrenciesRepo {
	return &currencies{
		conn: conn,
	}
}

func (c *currencies) GetList() ([]model.Currency, error) {
	var currencies = make([]model.Currency, 0)

	rows, err := c.conn.Query(selectCurrencyListSQL)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		currency := model.Currency{}
		err := rows.Scan(&currency.ID, currency.Symbol, currency.Name)
		if err != nil {
			return nil, err
		}

		currencies = append(currencies, currency)
	}

	return currencies, rows.Close()
}
