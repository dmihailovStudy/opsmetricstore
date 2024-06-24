package main

import (
	"github.com/dmihailovStudy/opsmetricstore/internal/config/server"
	"github.com/dmihailovStudy/opsmetricstore/internal/db"
	"github.com/dmihailovStudy/opsmetricstore/internal/handlers"
	"github.com/dmihailovStudy/opsmetricstore/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	var config server.Config

	// read envs
	config.Load()
	log.Info().Interface("config", config).Msg("main(): startup with config")

	// create empty storage
	memStorage := storage.CreateDefaultStorage()
	if config.Restore {
		localStorage, err := storage.ReadStorageFromFile(config.Path)
		if err != nil {
			log.Error().Err(err).Msg("main(): error while loading local snapshot")
		} else {
			memStorage = localStorage
		}
	}

	if config.SaveMode == "db" {
		db.ConnectPostgres(log.Logger, config.DBDSN)
		db.InitMigrations()
	}
	go storage.SaveStoragePeriodically(
		&memStorage,
		config.SaveMode,
		config.Path,
		time.Duration(config.StoreInterval)*time.Second,
	)

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.GET(server.MainPath, handlers.MainMiddleware(&memStorage))
	router.GET(server.GetMetricByURLPath, handlers.GetMetricByURLMiddleware(&memStorage))
	router.POST(server.GetMetricByJSONPath, handlers.GetMetricByJSONMiddleware(&memStorage))
	router.POST(server.UpdateByURLPath, handlers.UpdateByURLMiddleware(&memStorage))
	router.POST(server.UpdateByJSONPath, handlers.UpdateByJSONMiddleware(&memStorage))
	router.POST(server.UpdatesByJSONPath, handlers.UpdatesByJSONMiddleware(&memStorage))
	router.GET(server.GetDBStatusPath, handlers.GetDBStatusMiddleware())

	err := router.Run(config.Address)
	if err != nil {
		log.Error().Err(err).Msg("main(): router run error")
	}
}
