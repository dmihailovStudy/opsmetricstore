package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
)

type Con struct {
	Db *sql.DB
}

func ConnectPostgres(logger zerolog.Logger, dbDSN string) *Con {
	var con Con
	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		logger.Warn().Err(err).Msg("ConnectPostgres(): error while sql.Open")
	}
	con.Db = db
	defer db.Close()
	return &con
}

func (con *Con) Ping() error {
	return con.Db.Ping()
}
