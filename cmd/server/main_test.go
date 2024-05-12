package main

import (
	"fmt"
	"github.com/dmihailovStudy/opsmetricstore/internal/storage"
	"github.com/rs/zerolog/log"
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
	memStorage := storage.CreateDefaultStorage()

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

		code := storage.CheckUpdateMetricCorrectness(metricType, metricName, metricValue, &memStorage)
		testNum := i + 1
		if code != testCase.output {
			errLogHeader := fmt.Sprintf(
				"Test #%v (%s): failed - got: %v, want %v",
				testNum,
				testCase.name,
				code,
				testCase.output,
			)
			log.Error().Msg(errLogHeader)
			t.FailNow()
		}
		okLogHeader := fmt.Sprintf("Test #%v (%s): completed!\n", testNum, testCase.name)
		log.Info().Msg(okLogHeader)
	}
}
