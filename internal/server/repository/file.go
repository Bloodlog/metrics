package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type FileStorageWrapper struct {
	storage  MetricStorage
	filePath string
	interval int
}

func NewFileStorageWrapper(storage MetricStorage, filePath string, saveInterval int) *FileStorageWrapper {
	return &FileStorageWrapper{
		storage:  storage,
		filePath: filePath,
		interval: saveInterval,
	}
}

func (fw *FileStorageWrapper) SetGauge(ctx context.Context, name string, value float64) error {
	_ = fw.storage.SetGauge(ctx, name, value)
	if fw.interval > 0 {
		if err := fw.SaveToFile(ctx); err != nil {
			return fmt.Errorf("error saving metrics: %w", err)
		}
	}
	return nil
}

func (fw *FileStorageWrapper) GetGauge(ctx context.Context, name string) (float64, error) {
	value, err := fw.storage.GetGauge(ctx, name)
	if err != nil {
		return 0, fmt.Errorf("failed to get gauge '%s': %w", name, err)
	}
	return value, nil
}

func (fw *FileStorageWrapper) SetCounter(ctx context.Context, name string, value uint64) error {
	_ = fw.storage.SetCounter(ctx, name, value)
	if fw.interval > 0 {
		if err := fw.SaveToFile(ctx); err != nil {
			return fmt.Errorf("error saving counter: %w", err)
		}
	}

	return nil
}

func (fw *FileStorageWrapper) GetCounter(ctx context.Context, name string) (uint64, error) {
	value, err := fw.storage.GetCounter(ctx, name)
	if err != nil {
		return 0, fmt.Errorf("failed to get counter '%s': %w", name, err)
	}
	return value, nil
}

func (fw *FileStorageWrapper) Gauges(ctx context.Context) map[string]float64 {
	return fw.storage.Gauges(ctx)
}

func (fw *FileStorageWrapper) Counters(ctx context.Context) map[string]uint64 {
	return fw.storage.Counters(ctx)
}

func (fw *FileStorageWrapper) SaveToFile(ctx context.Context) error {
	data := struct {
		Gauges   map[string]float64 `json:"gauges"`
		Counters map[string]uint64  `json:"counters"`
	}{
		Gauges:   fw.storage.Gauges(ctx),
		Counters: fw.storage.Counters(ctx),
	}

	file, err := os.Create(fw.filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("error closing file %s: %w", fw.filePath, closeErr)
		}
	}()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode data to file: %w", err)
	}

	return nil
}

func (fw *FileStorageWrapper) LoadFromFile(ctx context.Context) error {
	file, err := os.Open(fw.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("error load file %s: %w", fw.filePath, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("error closing file %s: %w", fw.filePath, closeErr)
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
		err := fw.storage.SetGauge(ctx, k, v)
		if err != nil {
			return fmt.Errorf("error saving metrics: %w", err)
		}
	}
	for k, v := range data.Counters {
		err := fw.storage.SetCounter(ctx, k, v)
		if err != nil {
			return fmt.Errorf("error saving counter: %w", err)
		}
	}

	return nil
}

func (fw *FileStorageWrapper) AutoSave(ctx context.Context) error {
	if fw.interval <= 0 {
		return nil
	}

	ticker := time.NewTicker(time.Duration(fw.interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := fw.SaveToFile(ctx); err != nil {
			return fmt.Errorf("error saving: %w", err)
		}
	}

	return nil
}

func (fw *FileStorageWrapper) WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	return nil
}
