package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func ReadConfig(configFilePath string) (*Config, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	err = yaml.Unmarshal(data, cfg)
	return cfg, err
}
