package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // driver
)

const (
	driver         = "postgres"
	pssqlDsnFormat = "user=%s password=%s host=%s port=%s dbname=%s sslmode=disable"
)

// PostgresDSN connection model
type PostgresDSN struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

// Open opens connection with database
func Open(dsn *PostgresDSN) (*sql.DB, error) {
	db, err := sql.Open(driver, fmt.Sprintf(pssqlDsnFormat, dsn.User, dsn.Password, dsn.Server, dsn.Port, dsn.Database))
	if err != nil {
		return nil, err
	}

	// I want to be sure we have a real connection here
	return db, db.Ping()
}
