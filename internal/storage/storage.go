package storage

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"sync"
)

type Storage struct {
	Gauges   map[string]float64
	Counters map[string]int64
	gmx      *sync.RWMutex
	cmx      *sync.RWMutex
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

func GetMetricValue(metricType, metricName string, s *Storage) (bool, string, int64, float64, error) {
	err := errors.New("GetMetricValue: unknown metric type")
	metricValueString := ""
	isTracking := false
	metricValueInt := int64(0)
	metricValueFloat := float64(0)
	if metricType == CounterType {
		metricValueInt, isTracking = GetCounterMetric(metricName, s)
		metricValueString = fmt.Sprint(metricValueInt)
		err = nil
	} else if metricType == GaugeType {
		metricValueFloat, isTracking = GetGaugeMetric(metricName, s)
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

func CheckUpdateMetricCorrectness(metricType, metricName, metricValueStr string, s *Storage) int {
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
		storeValue, isTracking := GetCounterMetric(metricName, s)
		if !isTracking {
			UpdateCounterMetric(metricName, metricValueInt64, s)
		} else {
			UpdateCounterMetric(metricName, storeValue+metricValueInt64, s)
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
		s.Gauges[metricName] = metricValueFloat64
	} else {
		// bad metric type
		return http.StatusBadRequest
	}
	return http.StatusOK
}

func UpdateGaugeMetric(name string, value float64, s *Storage) {
	s.gmx.Lock()
	defer s.gmx.Unlock()
	s.Gauges[name] = value
}

func UpdateCounterMetric(name string, value int64, s *Storage) {
	s.cmx.Lock()
	defer s.cmx.Unlock()
	s.Counters[name] = value
}

func GetGaugeMetric(name string, s *Storage) (float64, bool) {
	s.gmx.RLock()
	defer s.gmx.Unlock()

	value, isTracking := s.Gauges[name]
	return value, isTracking
}

func GetCounterMetric(name string, s *Storage) (int64, bool) {
	s.cmx.RLock()
	defer s.cmx.Unlock()

	value, isTracking := s.Counters[name]
	return value, isTracking
}
