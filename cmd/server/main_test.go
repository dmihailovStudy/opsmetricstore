package main

import (
	"fmt"
	"github.com/dmihailovStudy/opsmetricstore/internal/metrics"
	"net/http"
	"strings"
	"testing"
)

type SendMetricTestCase struct {
	name   string
	input  string
	output int
}

func TestCheckMetricCorrectness(t *testing.T) {
	var memStorage metrics.Storage
	metrics.InitStorage(&memStorage)

	// Test #1: Ok case
	okTestName := "Ok case"
	okTestInput := "update/counter/metricName/0"
	okTestOutput := http.StatusOK
	okTest := SendMetricTestCase{
		okTestName,
		okTestInput,
		okTestOutput,
	}

	// Test #2: Incorrect metric type
	incorrectMetricTestName := "Incorrect metric type"
	incorrectMetricTestInput := "update/x/metricName/0"
	incorrectMetricTestOutput := http.StatusBadRequest
	incorrectMetricTypeTest := SendMetricTestCase{
		incorrectMetricTestName,
		incorrectMetricTestInput,
		incorrectMetricTestOutput,
	}

	// Test #3: Incorrect value case
	incorrectValueTestName := "Incorrect metric value type"
	incorrectValueTestInput := "update/counter/metricName/0.2"
	incorrectValueTestOutput := http.StatusBadRequest
	incorrectMetricValueTest := SendMetricTestCase{
		incorrectValueTestName,
		incorrectValueTestInput,
		incorrectValueTestOutput,
	}

	sendMetricTestCases := []SendMetricTestCase{
		okTest,
		incorrectMetricTypeTest,
		incorrectMetricValueTest,
	}

	for i, testCase := range sendMetricTestCases {
		metricData := strings.Split(testCase.input[1:], "/")
		_ = metricData[0]            // ex. "update"
		metricType := metricData[1]  // ex. "counter", "gauge"
		metricName := metricData[2]  // metricName to update
		metricValue := metricData[3] // metricValue in string format

		code := CheckUpdateMetricCorrectness(&memStorage, metricType, metricName, metricValue)
		if code != testCase.output {
			fmt.Printf("Test #%v (%s): failed - got: %v, want %v\n", i+1, testCase.name, code, testCase.output)
			t.FailNow()
		}
		fmt.Printf("Test #%v (%s): completed!\n", i+1, testCase.name)
	}
}
