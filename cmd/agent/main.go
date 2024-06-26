package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dmihailovStudy/opsmetricstore/internal/config/agent"
	"github.com/dmihailovStudy/opsmetricstore/internal/storage"
	"github.com/dmihailovStudy/opsmetricstore/transport/structure/metrics"
	"github.com/fatih/structs"
	"github.com/rs/zerolog/log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type UserStats struct {
	PollCount   int64
	RandomValue float64
}

const method = "update"
const compressRequest = true

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
		log.Err(err).Msg("main(): env load error")
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
				Msg("main(): new polling")
			runtime.ReadMemStats(&runtimeStats)
			userStats.PollCount = userStats.PollCount + 1
			userStats.RandomValue = rand.Float64()

		case reportTime := <-reportTicker.C:
			mapRuntimeStats := structs.Map(runtimeStats)
			mapUserStats := structs.Map(userStats)
			log.Info().
				Str("reportTime", reportTime.String()).
				Int64("pollCount", userStats.PollCount).
				Msg("main(): new reporting")
			sendMetrics(storage.RuntimeMetrics, mapRuntimeStats)
			sendMetrics(storage.UserMetrics, mapUserStats)
		}
	}
}

func sendMetrics(metricsArr []string, metricsMap map[string]interface{}) []string {
	var responsesStatus []string
	for _, metric := range metricsArr {
		metricType := storage.GetMetricType(metric)
		object := metrics.Body{
			ID:    metric,
			MType: metricType,
		}

		strMetric := fmt.Sprintf("%v", metricsMap[metric])

		if metricType == "counter" {
			value, err := strconv.ParseInt(strMetric, 10, 64)
			object.Delta = &value

			if err != nil {
				log.Error().Err(err).
					Str("name", metric).
					Str("value", strMetric).
					Msg("sendMetrics(): can't parse counter type")
			}

		} else {
			value, err := strconv.ParseFloat(strMetric, 64)
			if err != nil {
				log.Error().Err(err).
					Str("name", metric).
					Str("value", strMetric).
					Msg("sendMetrics(): can't parse gauge type")
			}

			object.Value = &value
		}

		objectBytes, err := json.Marshal(object)
		if err != nil {
			log.Error().Err(err).
				Msg("newConfig(): can't marshal new config")
		}

		var resp *http.Response
		contentType := "application/json"
		if compressRequest {
			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			if _, err := gz.Write(objectBytes); err != nil {
				strErr := fmt.Sprint(err)
				log.Error().
					Err(err).
					Str("path", baseURL).
					Msg("sendMetrics(): gz.Write error")
				responsesStatus = append(responsesStatus, strErr)
				continue
			}
			if err := gz.Close(); err != nil {
				strErr := fmt.Sprint(err)
				log.Error().
					Err(err).
					Str("path", baseURL).
					Msg("sendMetrics(): gz.Close() error")
				responsesStatus = append(responsesStatus, strErr)
				continue
			}

			// Отправка POST запроса с данными gzip на сервер
			req, err := http.NewRequest("POST", baseURL, &buf)
			if err != nil {
				log.Error().
					Err(err).
					Str("path", baseURL).
					Msg("sendMetrics(): build gzip request error")
			}
			req.Header.Set("Content-Encoding", "gzip")
			req.Header.Set("Content-Type", contentType)

			client := &http.Client{}
			resp, err = client.Do(req)
			if err != nil {
				strErr := fmt.Sprint(err)
				log.Error().
					Err(err).
					Str("path", baseURL).
					Msg("sendMetrics(): compressed response error")
				responsesStatus = append(responsesStatus, strErr)
				continue
			}
			defer resp.Body.Close()
		} else {
			body := bytes.NewBuffer(objectBytes)
			resp, err = http.Post(baseURL, contentType, body)
			if err != nil {
				strErr := fmt.Sprint(err)
				log.Error().
					Err(err).
					Str("path", baseURL).
					Msg("sendMetrics(): default post error")
				responsesStatus = append(responsesStatus, strErr)
				continue
			}
			defer resp.Body.Close()
		}

		log.Info().
			Str("path", baseURL).
			Str("status", resp.Status).
			Str("metricName", metric).
			Str("metricType", metricType).
			Msg("sendMetrics(): post ok")

		responsesStatus = append(responsesStatus, resp.Status)
	}
	return responsesStatus
}
