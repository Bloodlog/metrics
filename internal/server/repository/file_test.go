package repository

import (
	"os"
	"sync"
	"testing"
)

func TestSaveToFileAndLoadFromFile(t *testing.T) {
	const tempFileName = "metrics_test_*.json"
	const counterName = "NameCounter"
	const counterValue = 123

	const gaugeName = "NameGauge"
	const gaugeValue = 123.123
	memStorage := &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]uint64),
		mu:       &sync.RWMutex{},
	}

	err := memStorage.SetGauge(gaugeName, gaugeValue)
	if err != nil {
		t.Errorf("Failed to SetGauge: %v", err)
		return
	}
	err = memStorage.SetCounter(counterName, counterValue)
	if err != nil {
		t.Errorf("Failed to SetCounter: %v", err)
		return
	}

	tempFile, err := os.CreateTemp("", tempFileName)
	if err != nil {
		t.Errorf("Failed to create temp file: %v", err)
		return
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Errorf("Failed to remove temp file: %v", err)
			return
		}
	}(tempFile.Name())

	fileWrapper := NewFileStorageWrapper(memStorage, tempFile.Name(), 0)

	err = fileWrapper.SaveToFile()
	if err != nil {
		t.Errorf("SaveToFile failed: %v", err)
		return
	}

	memStorage = &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]uint64),
		mu:       &sync.RWMutex{},
	}
	fileWrapper.storage = memStorage

	err = fileWrapper.LoadFromFile()
	if err != nil {
		t.Errorf("LoadFromFile failed: %v", err)
		return
	}

	gauge, err := memStorage.GetGauge(gaugeName)
	if err != nil || gauge != gaugeValue {
		t.Errorf("Expected gauge value %v, got %v (err: %v)", gaugeValue, gauge, err)
	}

	counter, err := memStorage.GetCounter(counterName)
	if err != nil || counter != counterValue {
		t.Errorf("Expected counter value %v, got %v (err: %v)", counterValue, counter, err)
	}
}
