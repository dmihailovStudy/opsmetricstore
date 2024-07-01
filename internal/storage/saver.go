package storage

import (
	"encoding/json"
	"github.com/dmihailovStudy/opsmetricstore/internal/db"
	"github.com/dmihailovStudy/opsmetricstore/internal/db/models"
	"github.com/dmihailovStudy/opsmetricstore/internal/helpers"
	"github.com/dmihailovStudy/opsmetricstore/internal/retries"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
	"time"
)

func SaveStoragePeriodically(storage *Storage, saveMode, filePath string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		if saveMode == "db" {
			log.Info().Msg("SaveStoragePeriodically(): save snapshot to db")
			saveStorageToDB(storage)
		} else {
			log.Info().Msg("SaveStoragePeriodically(): write to snapshot file")
			saveStorageToJSONFile(storage, filePath)
		}
	}
}

func ReadStorageFromFile(filePath string) (Storage, error) {
	var storage Storage

	file, err := os.Open(filePath)
	if err != nil {
		return storage, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return storage, err
	}

	if err := json.Unmarshal(data, &storage); err != nil {
		return storage, err
	}

	return storage, nil
}

func saveStorageToJSONFile(storage *Storage, filePath string) {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		log.Warn().
			Err(errors.Unwrap(err)).
			Msg("saveStorageToJSONFile(): error creating directory")
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Warn().
			Str("path", filePath).
			Err(errors.Unwrap(err)).
			Msg("saveStorageToJSONFile(): error creating file")
		return
	}
	defer file.Close()

	data, err := json.Marshal(storage)
	if err != nil {
		log.Warn().
			Interface("storage", storage).
			Err(errors.Unwrap(err)).
			Msg("saveStorageToJSONFile(): error marshalling JSON")
		return
	}

	_, err = file.Write(data)
	if err != nil {
		log.Warn().
			Str("data", string(data)).
			Err(errors.Unwrap(err)).
			Msg("saveStorageToJSONFile(): error writing to file")
	}
}

func saveStorageToDB(storage *Storage) {
	go saveCountersToDB(storage.Counters)
	saveGaugesToDB(storage.Gauges)
}

func saveCountersToDB(storageGauges map[string]int64) {
	timestamp := time.Now().UTC()

	var counters models.Counters
	for name, value := range storageGauges {
		counter := models.Counter{Timestamp: timestamp, Name: name, Value: value}
		counters = append(counters, counter)
	}

	err := db.InsertCounters(counters)
	if err != nil {
		log.Warn().
			Interface("counters", counters).
			Err(err).
			Msg("saveCountersToDB(): failed to save counters")

		delayArr := []int{retries.FirstRetryDelay, retries.SecondRetryDelay, retries.ThirdRetryDelay}
		for i, delay := range delayArr {
			helpers.Wait(delay)
			err = db.InsertCounters(counters)
			if err != nil {
				log.Warn().
					Int("retry", i+1).
					Interface("counters", counters).
					Err(errors.Unwrap(err)).
					Msg("saveCountersToDB(): failed to retry counters")
			} else {
				break
			}
		}
	}
}

func saveGaugesToDB(storageGauges map[string]float64) {
	timestamp := time.Now().UTC()

	var gauges models.Gauges
	for name, value := range storageGauges {
		gauge := models.Gauge{Timestamp: timestamp, Name: name, Value: value}
		gauges = append(gauges, gauge)
	}

	err := db.InsertGauges(gauges)
	if err != nil {
		log.Warn().
			Interface("gauges", gauges).
			Err(errors.Unwrap(err)).
			Msg("saveGaugesToDB(): failed to save gauges")

		delayArr := []int{retries.FirstRetryDelay, retries.SecondRetryDelay, retries.ThirdRetryDelay}
		for i, delay := range delayArr {
			helpers.Wait(delay)
			err := db.InsertGauges(gauges)
			if err != nil {
				log.Warn().
					Int("retry", i+1).
					Interface("gauges", gauges).
					Err(errors.Unwrap(err)).
					Msg("saveGaugesToDB(): failed to retry gauges")
			} else {
				break
			}
		}
	}
}
