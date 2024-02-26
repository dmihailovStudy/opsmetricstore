package main

import (
	"flag"
	"github.com/dmihailovStudy/opsmetricstore/internal/config/server"
	"github.com/dmihailovStudy/opsmetricstore/internal/metrics"
	"github.com/dmihailovStudy/opsmetricstore/internal/templates/html"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

var memStorage metrics.Storage

func main() {
	var endpoint string
	flag.StringVar(&endpoint, server.AFlag, server.ADefault, server.AUsage)
	flag.Parse()

	var envs server.Envs
	err := envs.Load()
	if err != nil {
		log.Err(err).Msg("main: env load error")
	}

	if envs.Address != "" {
		endpoint = envs.Address
	}

	metrics.InitStorage(&memStorage)

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.SetHTMLTemplate(html.MetricsTemplate)
	router.GET(server.MainPath, MainHandler)
	router.GET(server.MetricPath, MetricHandler)
	router.POST(server.UpdatePath, UpdateContext)

	err = router.Run(endpoint)
	if err != nil {
		log.Err(err).Msg("main: router run error")
	}
}

func MainHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "metrics", gin.H{
		"gaugeBody":   memStorage.Gauge,
		"counterBody": memStorage.Counter,
	})
}

func MetricHandler(c *gin.Context) {
	metricType := c.Param("metricType")
	metricName := c.Param("metricName")

	isTracking, err, metricValueStr := metrics.GetMetricValueString(memStorage, metricType, metricName)
	if err != nil {
		log.Error().
			Err(err).
			Str("metricType", metricType).
			Str("metricName", metricType).
			Msg("Error: get metric value string")
	}

	log.Info().
		Bool("isTracking:", isTracking).
		Str("metricName", metricName).
		Str("metricType", metricType).
		Msg("New get metric request")

	if !isTracking {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	intCode, err := c.Writer.WriteString(metricValueStr)
	if err != nil {
		log.Error().
			Err(err).
			Str("metricValueStr", metricValueStr).
			Int("intCode", intCode).
			Msg("Error: while sending string")
	}
	c.Writer.WriteHeader(http.StatusOK)
}

func UpdateContext(c *gin.Context) {
	metricType := c.Param("metricType")
	metricName := c.Param("metricName")
	metricValue := c.Param("metricValue")

	responseCode := CheckUpdateMetricCorrectness(&memStorage, metricType, metricName, metricValue)
	c.Writer.WriteHeader(responseCode)
}

func CheckUpdateMetricCorrectness(memStorage *metrics.Storage, metricType, metricName, metricValueStr string) int {
	if metricType == metrics.CounterType {
		metricValueInt64, err := metrics.GetMetricValueInt64(metricValueStr)
		if err != nil {
			log.Error().
				Err(err).
				Str("metricValueStr", metricValueStr).
				Int("counterBase", metrics.CounterBase).
				Int("counterBitSize", metrics.CounterBitSize).
				Msg("GetMetricValueInt64: failed to convert metricValueStr")
			return http.StatusBadRequest
		}
		_, isTracking := memStorage.Counter[metricName]
		if !isTracking {
			memStorage.Counter[metricName] = metricValueInt64
		} else {
			memStorage.Counter[metricName] += metricValueInt64
		}
	} else if metricType == metrics.GaugeType {
		metricValueFloat64, err := metrics.GetMetricValueFloat64(metricValueStr)
		if err != nil {
			log.Error().
				Err(err).
				Str("metricValueStr", metricValueStr).
				Int("gaugeBitSize", metrics.GaugeBitSize).
				Msg("GetMetricValueFloat64: failed to convert metricValueStr")
			return http.StatusBadRequest
		}
		memStorage.Gauge[metricName] = metricValueFloat64
	} else {
		// bad metric type
		return http.StatusBadRequest
	}
	return http.StatusOK
}
