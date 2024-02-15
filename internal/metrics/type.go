package metrics

import "strconv"

// Counter adds new value to previous

const CounterType = "counter"
const CounterBase = 10
const CounterBitSize = 64

// Gauge replaces previous value

const GaugeType = "gauge"
const GaugeBitSize = 64

func CheckTypeAndValueCorrectness(metricType, metricValueStr string) bool {
	if metricType != CounterType && metricType != GaugeType {
		return false
	}

	if metricType == CounterType {
		_, err := strconv.ParseInt(metricValueStr, CounterBase, CounterBitSize) // Check: counter is int64 type
		if err != nil {
			return false
		}
	} else if metricType == GaugeType {
		_, err := strconv.ParseFloat(metricValueStr, GaugeBitSize) // Check: gauge is float64 type
		if err != nil {
			return false
		}
	}

	return true
}
