package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func NewDatabase(dsn string) (*DB, error) {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	return &DB{DB: db}, nil
}
