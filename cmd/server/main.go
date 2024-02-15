package main

import (
	"flag"
	"fmt"
	"github.com/dmihailovStudy/opsmetricstore/internal/metrics"
	"github.com/dmihailovStudy/opsmetricstore/internal/templates/html"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"strings"
)

var memStorage = make(map[string]string)

const aFlag = "a"

var endpoint string

func main() {
	flag.StringVar(&endpoint, aFlag, "localhost:8080", "specify the url")
	flag.Parse()

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.SetHTMLTemplate(html.MetricsTemplate)

	updatePagePath := "/update/:metricType/:metricName/:metricValue"
	router.POST(updatePagePath, UpdatePage)

	metricPagePath := "/value/:metricType/:metricName"
	router.GET(metricPagePath, MetricPage)

	mainPagePath := "/"
	router.GET(mainPagePath, MainPage)

	err := router.Run(endpoint)
	if err != nil {
		log.Err(err).Msg("router run error")
	}
}

func UpdatePage(c *gin.Context) {
	url := c.Request.URL
	responseCode := CheckUpdateMetricCorrectness(url.Path)
	c.Writer.WriteHeader(responseCode)
}

func MetricPage(c *gin.Context) {
	metricType := c.Param("metricType")
	metricName := c.Param("metricName")

	metricValue, isTracking := memStorage[metricName]
	log.Info().
		Bool("isTracking:", isTracking).
		Str("metricName", metricName).
		Str("metricType", metricType).
		Msg("New get metric request")

	if !isTracking {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	intCode, err := c.Writer.WriteString(metricValue)
	if err != nil {
		log.Error().
			Err(err).
			Str("metricValue", metricValue).
			Int("intCode", intCode).
			Msg("Error: while sending string")
	}
	c.Writer.WriteHeader(http.StatusOK)
}

func MainPage(c *gin.Context) {
	body := ""
	for metricName, metricValue := range memStorage { // Итерирует через карту (ключ, значение)
		body += fmt.Sprintf("%s:%s ", metricName, metricValue)
	}

	c.HTML(http.StatusOK, "metrics", gin.H{
		"metrics": body,
	})
}

func CheckUpdateMetricCorrectness(url string) int {
	url = url[1:] // Delete first "/" to escape metricData[0] = ""
	metricData := strings.Split(url[1:], "/")
	_ = metricData[0]            // ex. "update"
	metricType := metricData[1]  // ex. "counter", "gauge"
	metricName := metricData[2]  // metricName to update
	metricValue := metricData[3] // metricValue in string format

	if !metrics.CheckTypeAndValueCorrectness(metricType, metricValue) {
		return http.StatusBadRequest
	}

	if metricType == metrics.CounterType {
		_, isTracking := memStorage[metricName]
		if !isTracking {
			memStorage[metricName] = metricValue
		} else {
			previousValue, err := strconv.ParseInt(memStorage[metricName], metrics.CounterBase, metrics.CounterBitSize)
			if err != nil {
				log.Error().
					Err(err).
					Str("previousValue", memStorage[metricName]).
					Int("counterBase", metrics.CounterBase).
					Int("counterBitSize", metrics.CounterBitSize).
					Msg("CheckUpdateMetricCorrectness: failed to convert previousValue")
			}

			addValue, err := strconv.ParseInt(metricValue, metrics.CounterBase, metrics.CounterBitSize)
			if err != nil {
				log.Error().
					Err(err).
					Str("addValue", metricValue).
					Int("counterBase", metrics.CounterBase).
					Int("counterBitSize", metrics.CounterBitSize).
					Msg("CheckUpdateMetricCorrectness: failed to convert addValue")
			}
			memStorage[metricName] = fmt.Sprint(previousValue + addValue)
		}
	} else {
		memStorage[metricName] = metricValue
	}
	return http.StatusOK
}
