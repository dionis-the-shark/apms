package postgres

import (
	"database/sql"
	"fmt"

	apmsenv "github.com/dionis-the-shark/apms-env"
	_ "github.com/lib/pq"
)

const driverName = "postgres"

func New(cfg apmsenv.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBSSLMode,
	)
	db, err := sql.Open(driverName, connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
