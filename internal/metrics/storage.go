package metrics

import "strconv"

type Storage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

// Counter adds new value to previous

const CounterType = "counter"
const CounterBase = 10
const CounterBitSize = 64

// Gauge replaces previous value

const GaugeType = "gauge"
const GaugeBitSize = 64

func InitStorage(storage *Storage) {
	storage.Counter = make(map[string]int64)
	storage.Gauge = make(map[string]float64)
}

func GetMetricType(metric string) string {
	if metric == PollCountMetric {
		return "counter"
	}
	return "gauge"
}

func GetMetricValueString(storage Storage, metricType, metricName string) (bool, error, string) {
	metricValueString := ""
	metricValueInt := int64(0)
	metricValueFloat := float64(0)
	isTracking := false
	if metricType == CounterType {
		metricValueInt, isTracking = storage.Counter[metricName]
		metricValueString = strconv.FormatInt(metricValueInt, 16)
	} else if metricType == GaugeType {
		metricValueFloat, isTracking = storage.Gauge[metricName]
		metricValueString = strconv.FormatFloat(metricValueFloat, 'f', 2, 64)
	}

	return isTracking, nil, metricValueString
}

func GetMetricValueInt64(metricValueStr string) (int64, error) {
	metricValue, err := strconv.ParseInt(metricValueStr, CounterBase, CounterBitSize)
	return metricValue, err
}

func GetMetricValueFloat64(metricValueStr string) (float64, error) {
	metricValue, err := strconv.ParseFloat(metricValueStr, GaugeBitSize)
	return metricValue, err
}
