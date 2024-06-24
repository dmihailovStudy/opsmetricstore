package db

import (
	"embed"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func InitMigrations() {
	goose.SetBaseFS(embedMigrations)
	if err := goose.Up(conn, "migrations"); err != nil {
		log.Fatal().Err(err).Msg("initMigrations(): goose.Up error")
	}
	log.Info().Msg("initMigrations(): migrations applied")
}
