package repo

import "database/sql"

type Repository struct {
	CurrenciesRepo
	WalletsRepo
}

func NewRepository(conn *sql.DB) Repository {
	return Repository{
		CurrenciesRepo: NewCurrenciesRepo(conn),
		WalletsRepo:    NewWalletsRepo(conn),
	}
}
