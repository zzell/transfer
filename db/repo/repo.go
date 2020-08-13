package repo

import "database/sql"

// Repository app storage
type Repository struct {
	CurrenciesRepo
	WalletsRepo
}

// NewRepository constructor
func NewRepository(conn *sql.DB) Repository {
	return Repository{
		CurrenciesRepo: NewCurrenciesRepo(conn),
		WalletsRepo:    NewWalletsRepo(conn),
	}
}
