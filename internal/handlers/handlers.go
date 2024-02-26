package handlers

import (
	"github.com/dmihailovStudy/opsmetricstore/internal/metrics"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func MainMiddleware(storage *metrics.Storage) gin.HandlerFunc {
	return func(c *gin.Context) { MainHandler(c, storage) }
}

func MainHandler(c *gin.Context, storage *metrics.Storage) {
	c.HTML(http.StatusOK, "metrics", gin.H{
		"gaugeBody":   storage.Gauge,
		"counterBody": storage.Counter,
	})
}

func MetricMiddleware(storage *metrics.Storage) gin.HandlerFunc {
	return func(c *gin.Context) { MetricHandler(c, storage) }
}

func MetricHandler(c *gin.Context, storage *metrics.Storage) {
	metricType := c.Param("metricType")
	metricName := c.Param("metricName")

	isTracking, metricValueStr, err := metrics.GetMetricValueString(metricType, metricName, storage)
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

func UpdateMiddleware(storage *metrics.Storage) gin.HandlerFunc {
	return func(c *gin.Context) { UpdateHandler(c, storage) }
}

func UpdateHandler(c *gin.Context, storage *metrics.Storage) {
	metricType := c.Param("metricType")
	metricName := c.Param("metricName")
	metricValue := c.Param("metricValue")

	responseCode := metrics.CheckUpdateMetricCorrectness(metricType, metricName, metricValue, storage)
	c.Writer.WriteHeader(responseCode)
}
