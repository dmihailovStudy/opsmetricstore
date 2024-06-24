package storage

import (
	"encoding/json"
	"fmt"
	"github.com/dmihailovStudy/opsmetricstore/internal/db"
	"github.com/dmihailovStudy/opsmetricstore/internal/db/models"
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
		fmt.Println("Error creating directory:", err)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	data, err := json.Marshal(storage)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func saveStorageToDB(storage *Storage) {
	timestamp := time.Now().UTC()

	var gauges models.Gauges
	for name, value := range storage.Gauges {
		gauge := models.Gauge{Timestamp: timestamp, Name: name, Value: value}
		gauges = append(gauges, gauge)
	}

	var counters models.Counters
	for name, value := range storage.Counters {
		counter := models.Counter{Timestamp: timestamp, Name: name, Value: value}
		counters = append(counters, counter)
	}

	err := db.InsertGauges(gauges)
	if err != nil {
		log.Warn().Interface("gauges", gauges).Err(err).Msg("saveStorageToDB(): failed to save gauges")
	}

	err = db.InsertCounters(counters)
	if err != nil {
		log.Warn().Interface("counters", counters).Err(err).Msg("saveStorageToDB(): failed to save counters")
	}
}
