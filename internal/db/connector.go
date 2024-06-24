package db

import (
	"database/sql"
	"github.com/dmihailovStudy/opsmetricstore/internal/helpers"
	"github.com/dmihailovStudy/opsmetricstore/internal/retries"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var conn *sql.DB

func ConnectPostgres(log zerolog.Logger, dbDSN string) {
	var err error
	conn, err = sql.Open("pgx", dbDSN)
	if err != nil {
		log.Warn().Err(errors.Unwrap(err)).Msg("ConnectPostgres(): error while sql.Open")
	}

	err = conn.Ping()
	if err != nil {
		log.Warn().Err(errors.Unwrap(err)).Msg("ConnectPostgres(): error while sql.Ping")

		// retry ping
		delayArr := []int{retries.FirstRetryDelay, retries.SecondRetryDelay, retries.ThirdRetryDelay}
		for i, delay := range delayArr {
			helpers.Wait(delay)
			err = conn.Ping()
			if err != nil {
				log.Warn().
					Int("retry", i+1).
					Err(errors.Unwrap(err)).
					Msg("ConnectPostgres(): failed to retry sql.Ping")
			} else {
				break
			}
		}
	}
}

func Ping() error {
	return conn.Ping()
}
