package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dmihailovStudy/opsmetricstore/internal/config/agent"
	"github.com/dmihailovStudy/opsmetricstore/internal/metrics"
	"github.com/dmihailovStudy/opsmetricstore/internal/objects/update"
	"github.com/fatih/structs"
	"github.com/rs/zerolog/log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type UserStats struct {
	PollCount   int64
	RandomValue float64
}

const method = "update"

var baseURL string
var endpoint string
var pollIntervalSec int
var reportIntervalSec int

func main() {
	var envs agent.Envs
	var runtimeStats runtime.MemStats

	flag.StringVar(&endpoint, agent.AFlag, agent.ADefault, agent.AUsage)
	flag.IntVar(&pollIntervalSec, agent.PFlag, agent.PDefault, agent.PUsage)
	flag.IntVar(&reportIntervalSec, agent.RFlag, agent.RDefault, agent.RUsage)
	flag.Parse()

	err := envs.Load()
	if err != nil {
		log.Err(err).Msg("main: env load error")
	}

	if envs.Address != "" {
		endpoint = envs.Address
	}
	if envs.PollInterval != 0 {
		pollIntervalSec = envs.PollInterval
	}
	if envs.ReportInterval != 0 {
		reportIntervalSec = envs.ReportInterval
	}

	baseURL = fmt.Sprintf("http://%s/%s", endpoint, method)

	pollTicker := time.NewTicker(time.Duration(pollIntervalSec) * time.Second)
	reportTicker := time.NewTicker(time.Duration(reportIntervalSec) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	userStats := UserStats{0, rand.Float64()}
	for {
		select {
		case pollTime := <-pollTicker.C:
			log.Info().
				Str("pollTime", pollTime.String()).
				Msg("New polling")
			runtime.ReadMemStats(&runtimeStats)
			userStats.PollCount = userStats.PollCount + 1
			userStats.RandomValue = rand.Float64()

		case reportTime := <-reportTicker.C:
			mapRuntimeStats := structs.Map(runtimeStats)
			mapUserStats := structs.Map(userStats)
			log.Info().
				Str("reportTime", reportTime.String()).
				Int64("pollCount", userStats.PollCount).
				Msg("New reporting")
			sendMetrics(metrics.RuntimeMetrics, mapRuntimeStats)
			sendMetrics(metrics.UserMetrics, mapUserStats)
		}
	}
}

func sendMetrics(metricsArr []string, metricsMap map[string]interface{}) []string {
	var responsesStatus []string
	for _, metric := range metricsArr {
		metricType := metrics.GetMetricType(metric)
		object := update.MetricRequestObj{
			ID:    metric,
			MType: metricType,
		}

		if metricType == "counter" {
			object.Delta = fmt.Sprintf("%v", metricsMap[metric])
		} else {
			object.Value = fmt.Sprintf("%v", metricsMap[metric])
		}

		objectBytes, err := json.Marshal(object)
		if err != nil {
			log.Err(err).Msg("newConfig: can't marshal new config")
		}
		body := bytes.NewBuffer(objectBytes)

		resp, err := http.Post(baseURL, "text/plain", body)
		if err != nil {
			strErr := fmt.Sprint(err)
			log.Error().
				Err(err).
				Str("path", baseURL).
				Msg("sendMetrics: post error")
			responsesStatus = append(responsesStatus, strErr)
			continue
		}
		log.Info().
			Str("path", baseURL).
			Str("status", resp.Status).
			Msg("sendMetrics: post ok")
		responsesStatus = append(responsesStatus, resp.Status)
		defer resp.Body.Close()
	}
	return responsesStatus
}
