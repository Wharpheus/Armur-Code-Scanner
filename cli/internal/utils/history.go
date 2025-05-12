package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const historyFile = ".armur_history.json"

func getHistoryPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, historyFile)
}

func SaveToHistory(taskID string) error {
	history, _ := LoadHistory()
	history = append(history, taskID)
	data, _ := json.Marshal(history)
	return os.WriteFile(getHistoryPath(), data, 0644)
}

func LoadHistory() ([]string, error) {
	path := getHistoryPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []string{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var history []string
	json.Unmarshal(data, &history)
	return history, nil
}