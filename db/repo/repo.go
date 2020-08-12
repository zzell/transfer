package repo

import (
	"database/sql"
	"fmt"
)

const (
	selectByWalletID = `SELECT score FROM wallets WHERE id = $1`
	subtractScore    = `UPDATE wallets SET score = score - $1 WHERE id = $2`
	addScore         = `UPDATE wallets SET score = score + $1 WHERE id = $2`

	errRollbackFmt = "unable to rollback: %s, reason: %w"
)

type (
	// WalletsRepo represents access to wallets storage
	WalletsRepo interface {
		// instead of Transfer we could have "Add" and "Sub" methods but
		// we need to use single transaction to be sure that everything executed in one batch
		Transfer(sender, receiver int, score float64) error
	}

	wallets struct {
		conn *sql.DB
	}
)

// NewWalletsRepo constructor
func NewWalletsRepo(conn *sql.DB) WalletsRepo {
	return &wallets{conn: conn}
}

func (w *wallets) Transfer(sender, receiver int, score float64) error {
	tx, err := w.conn.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(subtractScore, score, sender)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return fmt.Errorf(errRollbackFmt, err2, err)
		}
		return err
	}

	_, err = tx.Exec(addScore, score, receiver)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return fmt.Errorf(errRollbackFmt, err2, err)
		}
		return err
	}

	return tx.Commit()
}

func (w *wallets) GetScore(walletID int) (float64, error) {
	var score float64

	err := w.conn.QueryRow(selectByWalletID, walletID).Scan(&score)
	if err != nil {
		return 0, err
	}

	return score, nil
}