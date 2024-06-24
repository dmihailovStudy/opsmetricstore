package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"testing"
)

type SendMetricsParams struct {
	metricsArr []string
	metricsMap map[string]interface{}
}

type SendMetricTestCase struct {
	name   string
	input  SendMetricsParams
	output []string
}

func TestSendMetrics(t *testing.T) {
	isServerRunning := false
	if !isServerRunning {
		t.Skipf("Server is dead")
	}

	// Test #1: Ok case
	name := "Ok case"
	metricsArr := []string{"Alloc"}
	metricsMap := map[string]interface{}{"Alloc": "4"}
	input := SendMetricsParams{metricsArr, metricsMap}
	output := []string{"200 OK"}
	okTest := SendMetricTestCase{name, input, output}

	// Test #2: Incorrect value case
	name = "Incorrect value case"
	metricsArr = []string{"Alloc"}
	metricsMap = map[string]interface{}{"Alloc": "hi"}
	input = SendMetricsParams{metricsArr, metricsMap}
	output = []string{"400 Bad Request"}
	incorrectTest := SendMetricTestCase{name, input, output}

	sendMetricTestCases := []SendMetricTestCase{
		okTest,
		incorrectTest,
	}

	for i, testCase := range sendMetricTestCases {
		metricResponses := sendSingleMetrics(testCase.input.metricsArr, testCase.input.metricsMap)
		testNum := i + 1
		for j, metricResponse := range metricResponses {
			if metricResponse != testCase.output[j] {
				errLogHeader := fmt.Sprintf(
					"Test #%v failed - got: %s, want %s\n",
					testNum,
					metricResponse,
					testCase.output[j],
				)
				log.Error().Msg(errLogHeader)
				t.FailNow()
			}
		}
		okLogHeader := fmt.Sprintf("Test #%v (%s)\n", i+1, testCase.name)
		log.Error().Msg(okLogHeader)
	}
}
