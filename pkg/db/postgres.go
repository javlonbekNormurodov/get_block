package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func NewPostgresConnection(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
