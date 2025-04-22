package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"metrics/internal/config"
	"os"
	"time"

	"go.uber.org/zap"
)

type FileStorageWrapper struct {
	storage MetricStorage
	cfg     *config.ServerConfig
	logger  *zap.SugaredLogger
}

func NewFileStorageWrapper(
	ctx context.Context,
	cfg *config.ServerConfig,
	logger *zap.SugaredLogger,
) (*FileStorageWrapper, error) {
	handlerLogger := logger.With("file", "NewFileStorageWrapper")
	memRepo, _ := NewMemStorage()

	fileStorage := &FileStorageWrapper{
		storage: memRepo,
		cfg:     cfg,
		logger:  handlerLogger,
	}
	if fileStorage.cfg.Restore {
		handlerLogger.Info("Restore is enabled, loading from file...")
		if err := fileStorage.loadFromFile(ctx); err != nil {
			return nil, &RetriableError{Err: err}
		}
		handlerLogger.Info("Successfully loaded metrics from file.")
	}

	if fileStorage.isEnableAutoSave() {
		go func() {
			handlerLogger.Info("autoSave enabled")
			if err := fileStorage.autoSave(ctx); err != nil {
				handlerLogger.Info("Error during autoSave", "error", err)
			}
		}()
	}

	handlerLogger.Infof("Using file storage: %s", fileStorage.cfg.FileStoragePath)
	return fileStorage, nil
}

func (fw *FileStorageWrapper) SetGauge(ctx context.Context, name string, value float64) (float64, error) {
	value, _ = fw.storage.SetGauge(ctx, name, value)
	if fw.isEnableAutoSave() {
		if err := fw.saveToFile(ctx); err != nil {
			return 0, fmt.Errorf("error set gauge: %w", err)
		}
	}
	return value, nil
}

func (fw *FileStorageWrapper) GetGauge(ctx context.Context, name string) (float64, error) {
	value, err := fw.storage.GetGauge(ctx, name)
	if err != nil {
		return 0, fmt.Errorf("failed to get gauge '%s': %w", name, err)
	}
	return value, nil
}

func (fw *FileStorageWrapper) SetCounter(ctx context.Context, name string, value uint64) (uint64, error) {
	value, _ = fw.storage.SetCounter(ctx, name, value)
	if fw.isEnableAutoSave() {
		if err := fw.saveToFile(ctx); err != nil {
			return 0, fmt.Errorf("error set counter: %w", err)
		}
	}

	return value, nil
}

func (fw *FileStorageWrapper) GetCounter(ctx context.Context, name string) (uint64, error) {
	value, err := fw.storage.GetCounter(ctx, name)
	if err != nil {
		return 0, fmt.Errorf("failed to get counter '%s': %w", name, err)
	}
	return value, nil
}

func (fw *FileStorageWrapper) Gauges(ctx context.Context) (map[string]float64, error) {
	gauges, err := fw.storage.Gauges(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gauges from storage: %w", err)
	}
	return gauges, nil
}

func (fw *FileStorageWrapper) Counters(ctx context.Context) (map[string]uint64, error) {
	counters, err := fw.storage.Counters(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get counters from storage: %w", err)
	}
	return counters, nil
}

func (fw *FileStorageWrapper) saveToFile(ctx context.Context) error {
	gauges, err := fw.storage.Gauges(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gauges: %w", err)
	}

	counters, err := fw.storage.Counters(ctx)
	if err != nil {
		return fmt.Errorf("failed to get counters: %w", err)
	}

	data := struct {
		Gauges   map[string]float64 `json:"gauges"`
		Counters map[string]uint64  `json:"counters"`
	}{
		Gauges:   gauges,
		Counters: counters,
	}

	file, err := os.Create(fw.cfg.FileStoragePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("error closing file %s: %w", fw.cfg.FileStoragePath, closeErr)
		}
	}()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode data to file: %w", err)
	}

	return nil
}

func (fw *FileStorageWrapper) UpdateCounterAndGauges(
	ctx context.Context,
	counters map[string]uint64,
	gauges map[string]float64,
) error {
	for counterName, counterValue := range counters {
		_, err := fw.storage.SetCounter(ctx, counterName, counterValue)
		if err != nil {
			return fmt.Errorf("error set counter: %w", err)
		}
	}

	for gaugeName, gaugeValue := range gauges {
		_, err := fw.storage.SetGauge(ctx, gaugeName, gaugeValue)
		if err != nil {
			return fmt.Errorf("error set gauge: %w", err)
		}
	}

	if fw.isEnableAutoSave() {
		if err := fw.saveToFile(ctx); err != nil {
			return fmt.Errorf("error save fail: %w", err)
		}
	}

	return nil
}

func (fw *FileStorageWrapper) isEnableAutoSave() bool {
	return fw.cfg.StoreInterval > 0
}

func (fw *FileStorageWrapper) loadFromFile(ctx context.Context) error {
	file, err := os.Open(fw.cfg.FileStoragePath)
	if err != nil {
		if os.IsNotExist(err) {
			fw.logger.Info("loadFromFile: File not exist.")
			return nil
		}

		return fmt.Errorf("error load file %s: %w", fw.cfg.FileStoragePath, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("error closing file %s: %w", fw.cfg.FileStoragePath, closeErr)
		}
	}()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", fw.cfg.FileStoragePath, err)
	}
	if info.Size() == 0 {
		fw.logger.Warn("loadFromFile: File is empty, skipping restore")
		return nil
	}

	var data struct {
		Gauges   map[string]float64 `json:"gauges"`
		Counters map[string]uint64  `json:"counters"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("failed to decode data from file: %w", err)
	}

	for k, v := range data.Gauges {
		_, err := fw.storage.SetGauge(ctx, k, v)
		if err != nil {
			return fmt.Errorf("error set gauges: %w", err)
		}
	}
	for k, v := range data.Counters {
		_, err := fw.storage.SetCounter(ctx, k, v)
		if err != nil {
			return fmt.Errorf("error set counters: %w", err)
		}
	}

	return nil
}

func (fw *FileStorageWrapper) autoSave(ctx context.Context) error {
	if fw.cfg.StoreInterval <= 0 {
		return nil
	}

	ticker := time.NewTicker(time.Duration(fw.cfg.StoreInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			_ = fw.saveToFile(context.Background())
			fw.logger.Info("AutoSave stopped due to context cancel")
			return nil
		case <-ticker.C:
			if err := fw.saveToFile(ctx); err != nil {
				return fmt.Errorf("error saving: %w", err)
			}
		}
	}
}
