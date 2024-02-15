package main

import (
	"fmt"
	"net/http"
	"testing"
)

type SendMetricTestCase struct {
	name   string
	input  string
	output int
}

func TestCheckMetricCorrectness(t *testing.T) {
	// Test #1: Ok case
	name := "Ok case"
	input := "update/counter/metricName/0"
	output := http.StatusOK
	okTest := SendMetricTestCase{name, input, output}

	// Test #2: Incorrect metric type
	name = "Incorrect metric type"
	input = "update/x/metricName/0"
	output = http.StatusBadRequest
	incorrectMetricTypeTest := SendMetricTestCase{name, input, output}

	// Test #3: Incorrect value case
	name = "Incorrect metric value type"
	input = "update/counter/metricName/0.2"
	output = http.StatusBadRequest
	incorrectMetricValueTest := SendMetricTestCase{name, input, output}

	sendMetricTestCases := []SendMetricTestCase{
		okTest,
		incorrectMetricTypeTest,
		incorrectMetricValueTest,
	}

	for i, testCase := range sendMetricTestCases {
		code := CheckUpdateMetricCorrectness(testCase.input)
		if code != testCase.output {
			fmt.Printf("Test #%v (%s): failed - got: %v, want %v\n", i+1, testCase.name, code, testCase.output)
			t.FailNow()
		}
		fmt.Printf("Test #%v (%s): completed!\n", i+1, testCase.name)
	}
}
