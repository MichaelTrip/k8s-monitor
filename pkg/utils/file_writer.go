package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type ChangeFileData struct {
	SavedAt time.Time `json:"savedAt"`
	Changes []interface{} `json:"changes"`
}

var (
	mu sync.Mutex
)

func SaveChangesToFile(filePath string, changes interface{}) error {
	mu.Lock()
	defer mu.Unlock()

	data := ChangeFileData{
		SavedAt: time.Now(),
		Changes: changes.([]interface{}),
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal changes: %v", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}

func LoadChangesFromFile(filePath string) ([]interface{}, error) {
	mu.Lock()
	defer mu.Unlock()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []interface{}{}, nil // Return empty slice if file doesn't exist
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var fileData ChangeFileData
	if err := json.Unmarshal(data, &fileData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal changes: %v", err)
	}

	return fileData.Changes, nil
}

func AppendChangeToFile(filePath string, change interface{}) error {
	mu.Lock()
	defer mu.Unlock()

	// Load existing changes
	var existingData ChangeFileData
	if data, err := os.ReadFile(filePath); err == nil {
		json.Unmarshal(data, &existingData)
	}

	// Append new change
	existingData.Changes = append(existingData.Changes, change)
	existingData.SavedAt = time.Now()

	// Keep only last 10000 changes to prevent file from growing too large
	if len(existingData.Changes) > 10000 {
		existingData.Changes = existingData.Changes[len(existingData.Changes)-10000:]
	}

	jsonData, err := json.MarshalIndent(existingData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal changes: %v", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}