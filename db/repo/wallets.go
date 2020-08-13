package repo

import (
	"database/sql"
	"fmt"

	"github.com/zzell/transfer/model"
)

const (
	subtractWalletScoreSQL = `UPDATE wallets SET score = score - $1 WHERE id = $2`
	addWalletScoreSQL      = `UPDATE wallets SET score = score + $1 WHERE id = $2`
	selectWalletByIDSQL    = `SELECT wallets.id, wallets.score, currencies.id, currencies.name, currencies.symbol
		FROM wallets
		LEFT JOIN currencies ON wallets.currency_id = currencies.id
		WHERE wallets.id = $1`

	errRollbackFmt = "unable to rollback: %s, reason: %w"
)

type (
	// WalletsRepo represents access to wallets storage
	WalletsRepo interface {
		// instead of Transfer we could have "Add" and "Sub" methods but
		// we need to use single transaction to be sure that everything executed in one batch
		Transfer(sender, receiver int, gross, net float64) error
		GetWallet(walletID int) (*model.Wallet, error)
	}

	wallets struct {
		conn *sql.DB
	}
)

// NewWalletsRepo constructor
func NewWalletsRepo(conn *sql.DB) WalletsRepo {
	return &wallets{conn: conn}
}

func (w *wallets) Transfer(sender, receiver int, sub, add float64) error {
	tx, err := w.conn.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(subtractWalletScoreSQL, sub, sender)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return fmt.Errorf(errRollbackFmt, err2, err)
		}
		return err
	}

	_, err = tx.Exec(addWalletScoreSQL, add, receiver)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return fmt.Errorf(errRollbackFmt, err2, err)
		}
		return err
	}

	return tx.Commit()
}

func (w *wallets) GetWallet(walletID int) (*model.Wallet, error) {
	var wallet = model.Wallet{}

	err := w.conn.QueryRow(selectWalletByIDSQL, walletID).Scan(&wallet.ID, &wallet.Score, &wallet.Currency.ID, &wallet.Currency.Name, &wallet.Currency.Symbol)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}
