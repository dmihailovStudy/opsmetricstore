package main

import (
	"flag"
	"github.com/dmihailovStudy/opsmetricstore/internal/config/server"
	"github.com/dmihailovStudy/opsmetricstore/internal/handlers"
	"github.com/dmihailovStudy/opsmetricstore/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	var endpoint string
	var interval int
	var path string
	var restore bool
	var envs server.Envs

	// read flags
	flag.StringVar(&endpoint, server.AFlag, server.ADefault, server.AUsage)
	flag.IntVar(&interval, server.IFlag, server.IDefault, server.IUsage)
	flag.StringVar(&path, server.FFlag, server.FDefault, server.FUsage)
	flag.BoolVar(&restore, server.RFlag, server.RDefault, server.RUsage)
	flag.Parse()

	// read envs
	err := envs.Load()
	if err != nil {
		log.Error().Err(err).Msg("main: env load error")
	}

	if envs.Address != "" {
		endpoint = envs.Address
	}
	if envs.StoreInterval != 0 {
		interval = envs.StoreInterval
	}
	if envs.Path != "" {
		path = envs.Path
	}
	if envs.Restore {
		restore = envs.Restore
	}

	// create empty storage
	memStorage := storage.CreateDefaultStorage()

	if restore {
		localStorage, err := storage.ReadStorageFromFile(path)
		if err != nil {
			log.Error().Err(err).Msg("main(): error while loading local snapshot")
		} else {
			memStorage = localStorage
		}
	}

	go storage.SaveStoragePeriodically(&memStorage, path, time.Duration(interval)*time.Second)

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.GET(server.MainPath, handlers.MainMiddleware(&memStorage))
	router.GET(server.GetMetricByURLPath, handlers.GetMetricByURLMiddleware(&memStorage))
	router.POST(server.GetMetricByJSONPath, handlers.GetMetricByJSONMiddleware(&memStorage))
	router.POST(server.UpdateByURLPath, handlers.UpdateByURLMiddleware(&memStorage))
	router.POST(server.UpdateByJSONPath, handlers.UpdateByJSONMiddleware(&memStorage))

	err = router.Run(endpoint)
	if err != nil {
		log.Error().Err(err).Msg("main(): router run error")
	}
}
