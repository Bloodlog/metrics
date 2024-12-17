package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type FileStorageWrapper struct {
	Storage  MetricStorage
	FilePath string
	Interval int
}

func NewFileStorageWrapper(storage MetricStorage, filePath string, saveInterval int) *FileStorageWrapper {
	return &FileStorageWrapper{
		Storage:  storage,
		FilePath: filePath,
		Interval: saveInterval,
	}
}

func (fw *FileStorageWrapper) SetGauge(name string, value float64) {
	fw.Storage.SetGauge(name, value)
	if fw.Interval > 0 {
		if err := fw.SaveToFile(); err != nil {
			fmt.Printf("Error saving metrics: %v\n", err)
		}
	}
}

func (fw *FileStorageWrapper) GetGauge(name string) (float64, error) {
	value, err := fw.Storage.GetGauge(name)
	if err != nil {
		return 0, fmt.Errorf("failed to get gauge '%s': %w", name, err)
	}
	return value, nil
}

func (fw *FileStorageWrapper) SetCounter(name string, value uint64) {
	fw.Storage.SetCounter(name, value)
	if fw.Interval > 0 {
		if err := fw.SaveToFile(); err != nil {
			fmt.Printf("Error saving metrics: %v\n", err)
		}
	}
}

func (fw *FileStorageWrapper) GetCounter(name string) (uint64, error) {
	value, err := fw.Storage.GetCounter(name)
	if err != nil {
		return 0, fmt.Errorf("failed to get counter '%s': %w", name, err)
	}
	return value, nil
}

func (fw *FileStorageWrapper) Gauges() map[string]float64 {
	return fw.Storage.Gauges()
}

func (fw *FileStorageWrapper) Counters() map[string]uint64 {
	return fw.Storage.Counters()
}

func (fw *FileStorageWrapper) SaveToFile() error {
	data := struct {
		Gauges   map[string]float64 `json:"gauges"`
		Counters map[string]uint64  `json:"counters"`
	}{
		Gauges:   fw.Storage.Gauges(),
		Counters: fw.Storage.Counters(),
	}

	file, err := os.Create(fw.FilePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode data to file: %w", err)
	}

	return nil
}

func (fw *FileStorageWrapper) LoadFromFile() error {
	file, err := os.Open(fw.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	var data struct {
		Gauges   map[string]float64 `json:"gauges"`
		Counters map[string]uint64  `json:"counters"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("failed to decode data from file: %w", err)
	}

	for k, v := range data.Gauges {
		fw.Storage.SetGauge(k, v)
	}
	for k, v := range data.Counters {
		fw.Storage.SetCounter(k, v)
	}

	return nil
}

func (fw *FileStorageWrapper) AutoSave() {
	if fw.Interval <= 0 {
		return
	}

	ticker := time.NewTicker(time.Duration(fw.Interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := fw.SaveToFile(); err != nil {
			fmt.Printf("Error saving metrics: %v\n", err)
		} else {
			fmt.Println("Metrics saved to file")
		}
	}
}
