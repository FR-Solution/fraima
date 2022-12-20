package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func ReadConfig(configFilePath string) ([]File, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var files []File
	err = yaml.Unmarshal(data, &files)
	return files, err
}
