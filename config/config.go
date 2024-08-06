package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type GitlabConfig struct {
	BaseUrl string `json:"baseUrl"`
	Token   string `json:"token"`
}

type Config struct {
	Gitlab GitlabConfig `json:"gitlab"`
}

const Filename = ".roehrich.json"

func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(homeDir, Filename)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s does not exist", configPath)
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
