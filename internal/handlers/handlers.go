package handlers

import (
	"encoding/json"
	"github.com/dmihailovStudy/opsmetricstore/internal/logging"
	"github.com/dmihailovStudy/opsmetricstore/internal/storage"
	"github.com/dmihailovStudy/opsmetricstore/transport/structure/metrics"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"time"
)

func MainMiddleware(s *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		lrw := MainHandler(c, s)
		lrw.LogQueryParams(c, startTime)
	}
}

func MainHandler(c *gin.Context, storage *storage.Storage) *logging.ResponseWriter {
	c.HTML(http.StatusOK, "storage", gin.H{
		"gaugeBody":   storage.Gauges,
		"counterBody": storage.Counters,
	})

	ginWriter := c.Writer
	lrw := logging.NewResponseWriter(ginWriter)
	return lrw
}

func GetMetricByJSONMiddleware(storage *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		lrw := GetMetricByJSONHandler(c, storage)
		lrw.LogQueryParams(c, startTime)
	}
}

func GetMetricByJSONHandler(c *gin.Context, s *storage.Storage) *logging.ResponseWriter {
	ginWriter := c.Writer
	loggingWriter := logging.NewResponseWriter(ginWriter)

	var requestObject metrics.Body
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("UpdateByJSONHandler(): io.ReadAll err")
	}

	err = json.Unmarshal(jsonData, &requestObject)
	if err != nil {
		log.Error().Err(err).Msg("UpdateByJSONHandler(): io.ReadAll err")
	}

	metricType := requestObject.MType
	metricName := requestObject.ID

	isTracking, metricValueStr, metricValueInt, metricValueFloat, err :=
		storage.GetMetricValue(metricType, metricName, s)

	if err != nil {
		log.Error().
			Err(err).
			Str("metricType", metricType).
			Str("metricName", metricType).
			Bool("isTracking", isTracking).
			Msg("Error: get metric get string")
	}

	log.Info().
		Str("metricName", metricName).
		Str("metricType", metricType).
		Msg("New get metric request")

	//if !isTracking {
	//	loggingWriter.WriteHeader(http.StatusNotFound)
	//	return loggingWriter
	//}

	var responseObject metrics.Body
	responseObject.ID = requestObject.ID
	responseObject.MType = requestObject.MType
	if metricType == "counter" {
		responseObject.Delta = &metricValueInt
	} else if metricType == "gauge" {
		responseObject.Value = &metricValueFloat
	}

	body, err := json.Marshal(responseObject)
	if err != nil {
		log.Error().
			Err(err).
			Str("metricName", metricName).
			Str("metricType", metricType).
			Bool("isTracking", isTracking).
			Str("metricValueStr", metricValueStr).
			Msg("Error: while marshal response")
	}

	loggingWriter.Header().Set("Content-Type", "application/json")
	_, err = loggingWriter.Write(body)
	if err != nil {
		log.Error().
			Err(err).
			Str("metricName", metricName).
			Str("metricType", metricType).
			Bool("isTracking", isTracking).
			Str("metricValueStr", metricValueStr).
			Msg("Error: while sending obj")
	}

	loggingWriter.WriteHeader(http.StatusOK)
	return loggingWriter
}

func GetMetricByURLMiddleware(storage *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		lrw := GetMetricByURLHandler(c, storage)
		lrw.LogQueryParams(c, startTime)
	}
}

func GetMetricByURLHandler(c *gin.Context, s *storage.Storage) *logging.ResponseWriter {
	ginWriter := c.Writer
	loggingWriter := logging.NewResponseWriter(ginWriter)

	metricType := c.Param("metricType")
	metricName := c.Param("metricName")

	isTracking, metricValueStr, _, _, err := storage.GetMetricValue(metricType, metricName, s)
	if err != nil {
		log.Error().
			Err(err).
			Str("metricType", metricType).
			Str("metricName", metricType).
			Msg("Error: get metric get string")
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

func UpdateByJSONMiddleware(s *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		lrw := UpdateByJSONHandler(c, s)
		lrw.LogQueryParams(c, startTime)
	}
}

func UpdateByJSONHandler(c *gin.Context, s *storage.Storage) *logging.ResponseWriter {
	ginWriter := c.Writer
	lrw := logging.NewResponseWriter(ginWriter)

	var requestObject metrics.Body
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Str("jsonData", string(jsonData)).Msg("UpdateByJSONHandler(): io.ReadAll err")
		lrw.WriteHeader(http.StatusNotFound)
		return lrw
	}

	err = json.Unmarshal(jsonData, &requestObject)
	if err != nil {
		log.Error().Err(err).Str("jsonData", string(jsonData)).Msg("UpdateByJSONHandler(): Unmarshal err")
		lrw.WriteHeader(http.StatusNotFound)
		return lrw
	}

	metricType := requestObject.MType
	metricName := requestObject.ID
	metricDelta := requestObject.Delta
	metricValue := requestObject.Value
	var responseCode int
	if metricType == "gauge" {
		storage.UpdateGaugeMetric(metricName, metricValue, s)
	} else if metricType == "counter" {
		storage.UpdateCounterMetric(metricName, metricDelta, s)
	} else {
		_, err = lrw.WriteString("Unknown metric type")
		if err != nil {
			log.Error().Err(err).Msg("UpdateByJSONHandler(): WriteString err")
		}
		lrw.WriteHeader(http.StatusNotFound)
		return lrw
	}

	lrw.WriteHeader(responseCode)
	return lrw
}

func UpdateByURLMiddleware(storage *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		lrw := UpdateByURLHandler(c, storage)
		lrw.LogQueryParams(c, startTime)
	}
}

func UpdateByURLHandler(c *gin.Context, s *storage.Storage) *logging.ResponseWriter {
	ginWriter := c.Writer
	lrw := logging.NewResponseWriter(ginWriter)
	metricType := c.Param("metricType")
	metricName := c.Param("metricName")
	metricValue := c.Param("metricValue")

	responseCode := storage.CheckUpdateMetricCorrectness(metricType, metricName, metricValue, s)
	lrw.WriteHeader(responseCode)
	return lrw
}
