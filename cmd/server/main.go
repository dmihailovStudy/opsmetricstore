package main

import (
	"flag"
	"github.com/dmihailovStudy/opsmetricstore/internal/config/server"
	"github.com/dmihailovStudy/opsmetricstore/internal/handlers"
	"github.com/dmihailovStudy/opsmetricstore/internal/metrics"
	"github.com/dmihailovStudy/opsmetricstore/internal/templates/html"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	var endpoint string
	var envs server.Envs

	// read flags
	flag.StringVar(&endpoint, server.AFlag, server.ADefault, server.AUsage)
	flag.Parse()

	// read envs
	err := envs.Load()
	if err != nil {
		log.Err(err).Msg("main: env load error")
	}
	if envs.Address != "" {
		endpoint = envs.Address
	}

	// create empty storage
	memStorage := metrics.CreateDefaultStorage()

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.SetHTMLTemplate(html.MetricsTemplate)
	router.GET(server.MainPath, handlers.MainMiddleware(&memStorage))
	router.GET(server.MetricPath, handlers.MetricMiddleware(&memStorage))
	router.POST(server.UpdatePath, handlers.UpdateMiddleware(&memStorage))

	err = router.Run(endpoint)
	if err != nil {
		log.Err(err).Msg("main: router run error")
	}
}
