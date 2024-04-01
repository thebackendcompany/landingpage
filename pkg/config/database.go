package config

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseConfig struct {
	ID  string
	URL string
	DSN string

	db *sqlx.DB
}

func ConnectDB(dsn, url, id string) (*DatabaseConfig, error) {
	db, err := sqlx.Connect(dsn, url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.DB.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	db.DB.SetConnMaxLifetime(10 * time.Minute)

	return &DatabaseConfig{
		ID:  id,
		URL: url,
		DSN: dsn,
		db:  db,
	}, nil
}

func (dc *DatabaseConfig) GetDB() *sqlx.DB {
	return dc.db
}
