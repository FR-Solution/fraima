package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func GetInstructionList(configFilePath string) ([]Instruction, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var cfg []Instruction
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}
