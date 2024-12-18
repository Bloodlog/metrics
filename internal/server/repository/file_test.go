package repository

import (
	"os"
	"sync"
	"testing"
)

func TestSaveToFileAndLoadFromFile(t *testing.T) {
	const CounterName = "NameCounter"
	const CounterValue = 123

	const GaugeName = "NameGauge"
	const GaugeValue = 123.123
	memStorage := &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]uint64),
		mu:       &sync.RWMutex{},
	}

	memStorage.SetGauge(GaugeName, GaugeValue)
	memStorage.SetCounter(CounterName, CounterValue)

	tempFile, err := os.CreateTemp("", "metrics_test_*.json")
	if err != nil {
		t.Errorf("Failed to create temp file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Errorf("Failed to remove temp file: %v", err)
		}
	}(tempFile.Name())

	fileWrapper := NewFileStorageWrapper(memStorage, tempFile.Name(), 0)

	err = fileWrapper.SaveToFile()
	if err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	memStorage = &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]uint64),
		mu:       &sync.RWMutex{},
	}
	fileWrapper.Storage = memStorage

	err = fileWrapper.LoadFromFile()
	if err != nil {
		t.Errorf("LoadFromFile failed: %v", err)
	}

	gauge, err := memStorage.GetGauge(GaugeName)
	if err != nil || gauge != GaugeValue {
		t.Errorf("Expected gauge value %v, got %v (err: %v)", GaugeValue, gauge, err)
	}

	counter, err := memStorage.GetCounter(CounterName)
	if err != nil || counter != CounterValue {
		t.Errorf("Expected counter value %v, got %v (err: %v)", CounterValue, counter, err)
	}
}
