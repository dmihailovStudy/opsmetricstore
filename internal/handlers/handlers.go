package handlers

import (
	"github.com/dmihailovStudy/opsmetricstore/internal/metrics"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func MainMiddleware(storage *metrics.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		MainHandler(c, storage)
		logQueryParams(c, startTime)
	}
}

func MainHandler(c *gin.Context, storage *metrics.Storage) {
	c.HTML(http.StatusOK, "metrics", gin.H{
		"gaugeBody":   storage.Gauges,
		"counterBody": storage.Counters,
	})
}

func MetricMiddleware(storage *metrics.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		MetricHandler(c, storage)
		logQueryParams(c, startTime)
	}
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
	return func(c *gin.Context) {
		startTime := time.Now()
		UpdateHandler(c, storage)
		logQueryParams(c, startTime)
	}
}

func UpdateHandler(c *gin.Context, storage *metrics.Storage) {
	metricType := c.Param("metricType")
	metricName := c.Param("metricName")
	metricValue := c.Param("metricValue")

	responseCode := metrics.CheckUpdateMetricCorrectness(metricType, metricName, metricValue, storage)
	c.Writer.WriteHeader(responseCode)
}

func logQueryParams(c *gin.Context, startTime time.Time) {
	uri := c.Request.RequestURI
	method := c.Request.Method
	log.Info().
		Str("uri", uri).
		Str("method", method).
		Dur("execTime", time.Since(startTime)).
		Msg("logQueryParams(): req stats")
}
