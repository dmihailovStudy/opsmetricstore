package main

import (
	"bytes"
	"fmt"
	"github.com/dmihailovStudy/opsmetricstore/internal/metrics"
	"github.com/fatih/structs"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type UserStats struct {
	PollCount   int64
	RandomValue float64
}

const pollIntervalSec = 2    // Update metrics interval
const reportIntervalSec = 10 // Send metrics interval

const host = "localhost"
const port = "8080"
const method = "update"

var baseURL = fmt.Sprintf("http://%s:%s/%s", host, port, method)

func main() {
	poolTicker := time.NewTicker(pollIntervalSec * time.Second)
	reportTicker := time.NewTicker(reportIntervalSec * time.Second)
	defer poolTicker.Stop()

	var runtimeStats runtime.MemStats
	userStats := UserStats{0, rand.Float64()}

	for {
		select {
		case pollTime := <-poolTicker.C:
			fmt.Println("New pooling: ", pollTime)
			runtime.ReadMemStats(&runtimeStats)
			userStats.PollCount = userStats.PollCount + 1
			userStats.RandomValue = rand.Float64()

		case reportTime := <-reportTicker.C:
			mapRuntimeStats := structs.Map(runtimeStats)
			mapUserStats := structs.Map(userStats)
			fmt.Println("New reporting: ", reportTime, userStats.PollCount)
			sendMetrics(metrics.RuntimeMetrics, mapRuntimeStats)
			sendMetrics(metrics.UserMetrics, mapUserStats)
		}
	}
}

func sendMetrics(metricsArr []string, metricsMap map[string]interface{}) []string {
	var responsesStatus []string
	for _, metric := range metricsArr {
		metricType := metrics.GetMetricType(metric)
		path := fmt.Sprintf("%s/%s/%s/%v", baseURL, metricType, metric, metricsMap[metric])
		body := bytes.NewBuffer([]byte{})
		resp, err := http.Post(path, "text/plain", body)
		if err != nil {
			strErr := fmt.Sprint(err)
			fmt.Printf("%s, err: %s", path, strErr)
			responsesStatus = append(responsesStatus, strErr)
			continue
		}
		defer resp.Body.Close()
		fmt.Printf("%s, status: %s\n", path, resp.Status)
		responsesStatus = append(responsesStatus, resp.Status)
	}
	return responsesStatus
}
