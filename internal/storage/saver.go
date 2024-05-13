package storage

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func SaveStoragePeriodically(storage *Storage, filePath string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			log.Info().Msg("SaveStoragePeriodically(): write to snapshot file")
			saveStorageToJsonFile(storage, filePath)
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

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return storage, err
	}

	if err := json.Unmarshal(data, &storage); err != nil {
		return storage, err
	}

	return storage, nil
}

func saveStorageToJsonFile(storage *Storage, filePath string) {
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
