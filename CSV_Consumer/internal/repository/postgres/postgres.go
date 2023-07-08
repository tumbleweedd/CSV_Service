package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	PgPort    string
	PgHost    string
	PgDBName  string
	PgUser    string
	PgPwd     string
	PgSSLMode string
}

func NewPostgresDB(cfg *Config) (*sqlx.DB, error) {
	db, err := sqlx.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.PgHost, cfg.PgPort, cfg.PgUser, cfg.PgDBName, cfg.PgPwd, cfg.PgSSLMode),
	)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}
