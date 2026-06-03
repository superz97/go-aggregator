package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(*c)
}

func getConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(path, configFileName)
	return fullPath, nil
}

func Read() (Config, error) {
	config, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(config)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	configuration := Config{}
	err = decoder.Decode(&configuration)
	if err != nil {
		return Config{}, err
	}

	return configuration, nil
}

func write(cfg Config) error {
	config, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(config)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}
	return nil
}
