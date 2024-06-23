package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/dmihailovStudy/opsmetricstore/internal/logging"
	"github.com/dmihailovStudy/opsmetricstore/internal/storage"
	"github.com/dmihailovStudy/opsmetricstore/internal/templates/html"
	"github.com/dmihailovStudy/opsmetricstore/transport/structure/metrics"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func MainMiddleware(s *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		lrw := logging.NewResponseWriter(c.Writer)
		MainHandler(c, lrw, s)
		lrw.LogQueryParams(c, startTime)
	}
}

func MainHandler(c *gin.Context, lrw *logging.ResponseWriter, storage *storage.Storage) {
	acceptEncoding := c.Request.Header.Get("Accept-Encoding")
	t := template.Must(template.New("storage").Parse(html.MetricsTemplate))

	log.Info().
		Str("acceptEncoding", acceptEncoding).
		Interface("storage", storage).
		Msg("MainHandler(): log params")

	if strings.Contains(acceptEncoding, "gzip") {
		var buf bytes.Buffer
		err := t.Execute(&buf, &storage)
		if err != nil {
			log.Error().Err(err).Msg("MainHandler(): error while html/template gzip execute")
		}

		bytesResponse, err := EncodeResponse(buf.Bytes())
		if err != nil {
			log.Error().Err(err).Msg("MainHandler(): error while encodeResponse")
		}

		log.Info().
			Str("buf", buf.String()).
			Str("bytesResponse", string(bytesResponse)).
			Msg("MainHandler(): log buffer")

		lrw.SendEncodedBody(http.StatusOK, "text/html", bytesResponse)
	} else {
		err := t.Execute(lrw, &storage)
		if err != nil {
			log.Error().Err(err).Msg("MainHandler(): error while html/template gzip execute")
		}
	}
}

func GetMetricByJSONMiddleware(s *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		lrw := logging.NewResponseWriter(c.Writer)
		body, err := DecodeBody(c)

		if err != nil {
			log.Error().
				Err(err).
				Interface("body", body).
				Msg("UpdateByJSONMiddleware(): DecodeBody err")
			lrw.WriteHeader(http.StatusNotFound)
		} else {
			status, rawResponse := GetMetricByJSONHandler(body, s)
			PrepareAndSendResponse(c, lrw, status, rawResponse)
		}

		lrw.LogQueryParams(c, startTime)
	}
}

func GetMetricByJSONHandler(requestObject metrics.Body, s *storage.Storage) (int, []byte) {
	metricType := requestObject.MType
	metricName := requestObject.ID

	isTracking, metricValueStr, metricValueInt, metricValueFloat, err :=
		storage.GetMetricValue(metricType, metricName, s)

	if err != nil {
		errMsg := "GetMetricByJSONHandler(): get metric value error"
		log.Error().
			Err(err).
			Str("metricType", metricType).
			Str("metricName", metricType).
			Bool("isTracking", isTracking).
			Msg(errMsg)
		return http.StatusBadRequest, []byte(errMsg)
	}

	log.Info().
		Str("metricName", metricName).
		Str("metricType", metricType).
		Msg("GetMetricByJSONHandler(): new get metric request")

	if !isTracking {
		return http.StatusNotFound, []byte("GetMetricByJSONHandler(): not tracking metric")
	}

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
		errMsg := "GetMetricByJSONHandler(): err while marshal response"
		log.Error().
			Err(err).
			Str("metricName", metricName).
			Str("metricType", metricType).
			Bool("isTracking", isTracking).
			Str("metricValueStr", metricValueStr).
			Msg(errMsg)
		return http.StatusBadRequest, []byte(errMsg)
	}

	return http.StatusOK, body
}

func GetMetricByURLMiddleware(storage *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		lrw := logging.NewResponseWriter(c.Writer)
		status, rawResponse := GetMetricByURLHandler(c, storage)
		PrepareAndSendResponse(c, lrw, status, rawResponse)
		lrw.LogQueryParams(c, startTime)
	}
}

func GetMetricByURLHandler(c *gin.Context, s *storage.Storage) (int, []byte) {
	metricType := c.Param("metricType")
	metricName := c.Param("metricName")

	isTracking, metricValueStr, _, _, err := storage.GetMetricValue(metricType, metricName, s)
	if err != nil {
		errMsg := "GetMetricByURLHandler(): get metric get string"
		log.Error().
			Err(err).
			Str("metricType", metricType).
			Str("metricName", metricType).
			Msg(errMsg)
		return http.StatusBadRequest, []byte(errMsg)
	}

	if !isTracking {
		return http.StatusNotFound, []byte("GetMetricByURLHandler(): not tracking metric")
	}

	log.Info().
		Bool("isTracking:", isTracking).
		Str("metricName", metricName).
		Str("metricType", metricType).
		Msg("GetMetricByURLHandler(): new get metric request")

	return http.StatusOK, []byte(metricValueStr)
}

func UpdateByJSONMiddleware(s *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		lrw := logging.NewResponseWriter(c.Writer)
		body, err := DecodeBody(c)

		if err != nil {
			log.Error().
				Err(err).
				Interface("body", body).
				Msg("UpdateByJSONMiddleware(): DecodeBody err")
			lrw.WriteHeader(http.StatusNotFound)
		} else {
			status, rawResponse := UpdateByJSONHandler(body, s)
			PrepareAndSendResponse(c, lrw, status, rawResponse)
		}

		lrw.LogQueryParams(c, startTime)
	}
}

func UpdateByJSONHandler(requestObject metrics.Body, s *storage.Storage) (int, []byte) {
	responseStr := ""
	metricType := requestObject.MType
	metricName := requestObject.ID
	metricDelta := requestObject.Delta
	metricValue := requestObject.Value
	if metricType == "gauge" {
		storage.UpdateGaugeMetric(metricName, metricValue, s)
	} else if metricType == "counter" {
		storage.UpdateCounterMetric(metricName, metricDelta, s)
	} else {
		responseStr = "Unknown metric type"
		return http.StatusNotFound, []byte(responseStr)
	}
	return http.StatusOK, []byte(responseStr)
}

func UpdateByURLMiddleware(storage *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		lrw := logging.NewResponseWriter(c.Writer)
		status, rawResponse := UpdateByURLHandler(c, storage)
		PrepareAndSendResponse(c, lrw, status, rawResponse)
		lrw.LogQueryParams(c, startTime)
	}
}

func UpdateByURLHandler(c *gin.Context, s *storage.Storage) (int, []byte) {
	metricType := c.Param("metricType")
	metricName := c.Param("metricName")
	metricValue := c.Param("metricValue")

	status := storage.CheckUpdateMetricCorrectness(metricType, metricName, metricValue, s)

	return status, []byte(strconv.Itoa(status))
}

func DecodeBody(c *gin.Context) (metrics.Body, error) {
	var requestObject metrics.Body
	var err error
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().
			Err(err).
			Str("jsonData", string(jsonData)).
			Msg("DecodeBody(): io.ReadAll err")
	}

	encoding := c.Request.Header.Get("Content-Encoding")
	if encoding == "" {
		err = json.Unmarshal(jsonData, &requestObject)
		if err != nil {
			log.Error().
				Err(err).
				Str("jsonData", string(jsonData)).
				Msg("DecodeBody(): Unmarshal jsonData err")
		}
	} else if encoding == "gzip" {
		reader, err := gzip.NewReader(bytes.NewBuffer(jsonData))
		if err != nil {
			log.Error().
				Err(err).
				Str("jsonData", string(jsonData)).
				Msg("DecodeBody(): gzip.NewReader err")
		}

		if err := json.NewDecoder(reader).Decode(&requestObject); err != nil {
			log.Error().
				Err(err).
				Str("jsonData", string(jsonData)).
				Msg("DecodeBody(): decode err")
		}

		defer reader.Close()
	}

	log.Info().
		Interface("body", requestObject).
		Msg("DecodeBody(): body after decoding")

	return requestObject, err
}

func EncodeResponse(response []byte) ([]byte, error) {
	// Сжатие данных в формат gzip
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write(response)
	if err != nil {
		log.Error().Err(err).Msg("EncodeResponse(): error while gz.Write")
		return nil, err
	}
	if err := gz.Close(); err != nil {
		log.Error().Err(err).Msg("EncodeResponse(): error while gz.Close")
		return nil, err
	}

	return buf.Bytes(), nil
}

func PrepareAndSendResponse(c *gin.Context, lrw *logging.ResponseWriter, status int, rawResponse []byte) {
	acceptEncoding := c.Request.Header.Get("Accept-Encoding")
	log.Info().
		Str("acceptEncoding", acceptEncoding).
		Int("status", status).
		Bytes("rawResponse", rawResponse).
		Msg("PrepareAndSendResponse(): log input params")

	if string(rawResponse) == "" {
		lrw.WriteHeader(status)
	} else if strings.Contains(acceptEncoding, "gzip") {
		bytesResponse, err := EncodeResponse(rawResponse)
		if err != nil {
			lrw.WriteHeader(http.StatusBadRequest)
			intCode, err := lrw.WriteString(err.Error())
			if err != nil {
				log.Error().
					Err(err).
					Int("intCode", intCode).
					Msg("PrepareAndSendResponse(): Writer.WriteString error")
			}
		} else {
			lrw.SendEncodedBody(status, "application/json", bytesResponse)
		}
	} else {
		lrw.WriteHeader(status)
		_, err := lrw.Write(rawResponse)
		if err != nil {
			log.Error().
				Err(err).
				Msg("PrepareAndSendResponse(): Writer.Write error")
		}
	}
}
