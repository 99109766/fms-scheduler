package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	NumResources int     `yaml:"num_resources"`
	NumTasks     int     `yaml:"num_tasks"`
	TotalUtil    float64 `yaml:"total_util"`
}

// LoadConfig reads the YAML configuration file from the given path and returns a pointer to a Config struct.
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
