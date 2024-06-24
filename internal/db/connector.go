package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
)

var conn *sql.DB

func ConnectPostgres(log zerolog.Logger, dbDSN string) {
	var err error
	conn, err = sql.Open("pgx", dbDSN)
	if err != nil {
		log.Warn().Err(err).Msg("ConnectPostgres(): error while sql.Open")
	}

	err = conn.Ping()
	if err != nil {
		log.Warn().Err(err).Msg("ConnectPostgres(): error while sql.Ping")
	}
}

func Ping() error {
	return conn.Ping()
}
