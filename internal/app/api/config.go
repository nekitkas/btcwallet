package api

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port           string `json:"port"`
	DatabaseURL    string `json:"database_url"`
	DatabaseDriver string `json:"database_driver"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) ReadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return nil
}
