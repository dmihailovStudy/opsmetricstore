package metrics

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

type Storage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

// Counters adds new value to previous

const CounterType = "counter"
const CounterBase = 10
const CounterBitSize = 64

// Gauges replaces previous value

const GaugeType = "gauge"
const GaugeBitSize = 64

func CreateDefaultStorage() Storage {
	var storage Storage
	storage.Counters = make(map[string]int64)
	storage.Gauges = make(map[string]float64)
	return storage
}

func GetMetricType(metric string) string {
	if metric == PollCountMetric {
		return "counter"
	}
	return "gauge"
}

func GetMetricValue(metricType, metricName string, storage *Storage) (bool, string, int64, float64, error) {
	err := errors.New("GetMetricValue: unknown metric type")
	metricValueString := ""
	isTracking := false
	metricValueInt := int64(0)
	metricValueFloat := float64(0)
	if metricType == CounterType {
		metricValueInt, isTracking = storage.Counters[metricName]
		metricValueString = fmt.Sprint(metricValueInt)
		err = nil
	} else if metricType == GaugeType {
		metricValueFloat, isTracking = storage.Gauges[metricName]
		metricValueString = fmt.Sprint(metricValueFloat)
		err = nil
	}
	return isTracking, metricValueString, metricValueInt, metricValueFloat, err
}

func GetMetricValueInt64(metricValueStr string) (int64, error) {
	metricValue, err := strconv.ParseInt(metricValueStr, CounterBase, CounterBitSize)
	return metricValue, err
}

func GetMetricValueFloat64(metricValueStr string) (float64, error) {
	metricValue, err := strconv.ParseFloat(metricValueStr, GaugeBitSize)
	return metricValue, err
}

func CheckUpdateMetricCorrectness(metricType, metricName, metricValueStr string, storage *Storage) int {
	if metricType == CounterType {
		metricValueInt64, err := GetMetricValueInt64(metricValueStr)
		if err != nil {
			log.Error().
				Err(err).
				Str("metricValueStr", metricValueStr).
				Int("counterBase", CounterBase).
				Int("counterBitSize", CounterBitSize).
				Msg("GetMetricValueInt64: failed to convert metricValueStr")
			return http.StatusBadRequest
		}
		_, isTracking := storage.Counters[metricName]
		if !isTracking {
			storage.Counters[metricName] = metricValueInt64
		} else {
			storage.Counters[metricName] += metricValueInt64
		}
	} else if metricType == GaugeType {
		metricValueFloat64, err := GetMetricValueFloat64(metricValueStr)
		if err != nil {
			log.Error().
				Err(err).
				Str("metricValueStr", metricValueStr).
				Int("gaugeBitSize", GaugeBitSize).
				Msg("GetMetricValueFloat64: failed to convert metricValueStr")
			return http.StatusBadRequest
		}
		storage.Gauges[metricName] = metricValueFloat64
	} else {
		// bad metric type
		return http.StatusBadRequest
	}
	return http.StatusOK
}
