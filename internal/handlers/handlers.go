package handlers

import (
	"github.com/dmihailovStudy/opsmetricstore/internal/logging"
	"github.com/dmihailovStudy/opsmetricstore/internal/metrics"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func MainMiddleware(storage *metrics.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		lrw := MainHandler(c, storage)
		lrw.LogQueryParams(c, startTime)
	}
}

func MainHandler(c *gin.Context, storage *metrics.Storage) *logging.ResponseWriter {
	c.HTML(http.StatusOK, "metrics", gin.H{
		"gaugeBody":   storage.Gauges,
		"counterBody": storage.Counters,
	})

	ginWriter := c.Writer
	lrw := logging.NewResponseWriter(ginWriter)
	return lrw
}

func MetricMiddleware(storage *metrics.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		lrw := MetricHandler(c, storage)
		lrw.LogQueryParams(c, startTime)
	}
}

func MetricHandler(c *gin.Context, storage *metrics.Storage) *logging.ResponseWriter {
	ginWriter := c.Writer
	loggingWriter := logging.NewResponseWriter(ginWriter)

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
		loggingWriter.WriteHeader(http.StatusNotFound)
		return loggingWriter
	}

	intCode, err := loggingWriter.WriteString(metricValueStr)
	if err != nil {
		log.Error().
			Err(err).
			Str("metricValueStr", metricValueStr).
			Int("intCode", intCode).
			Msg("Error: while sending string")
	}

	loggingWriter.WriteHeader(http.StatusOK)
	return loggingWriter
}

func UpdateMiddleware(storage *metrics.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		lrw := UpdateHandler(c, storage)
		lrw.LogQueryParams(c, startTime)
	}
}

func UpdateHandler(c *gin.Context, storage *metrics.Storage) *logging.ResponseWriter {
	ginWriter := c.Writer
	lrw := logging.NewResponseWriter(ginWriter)
	metricType := c.Param("metricType")
	metricName := c.Param("metricName")
	metricValue := c.Param("metricValue")

	responseCode := metrics.CheckUpdateMetricCorrectness(metricType, metricName, metricValue, storage)
	lrw.WriteHeader(responseCode)
	return lrw
}
