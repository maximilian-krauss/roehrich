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
	Credentials map[string]GitlabConfig `json:"credentials"`
}

const Filename = ".roehrich.json"

func GetConfigByHostname(hostname string, config Config) (*GitlabConfig, error) {
	credential, ok := config.Credentials[hostname]
	if !ok {
		return nil, fmt.Errorf("no credential for hostname %s", hostname)
	}
	return &credential, nil
}

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
