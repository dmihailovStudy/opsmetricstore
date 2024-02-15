package main

import (
	"github.com/dmihailovStudy/opsmetricstore/internal/metrics"
	"net/http"
	"strings"
)

func updatePage(res http.ResponseWriter, req *http.Request) {
	isMethodCorrect := req.Method == http.MethodPost
	if !isMethodCorrect {
		res.WriteHeader(http.StatusNotFound)
	}

	url := req.URL
	responseCode := checkMetricCorrectness(url.Path)
	res.WriteHeader(responseCode)
}

func checkMetricCorrectness(url string) int {
	url = url[1:] // Delete first "/" to escape metricData[0] = ""
	metricData := strings.Split(url[1:], "/")
	if len(metricData) != metrics.ParamsNumber {
		return http.StatusNotFound
	}
	_ = metricData[0]            // ex. "update"
	metricType := metricData[1]  // ex. "counter", "gauge"
	_ = metricData[2]            // metricName to update
	metricValue := metricData[3] // metricValue in string format

	if !metrics.CheckTypeAndValueCorrectness(metricType, metricValue) {
		return http.StatusBadRequest
	}

	return http.StatusOK
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, updatePage)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
