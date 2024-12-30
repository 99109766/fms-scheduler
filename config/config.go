package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	NumResources  int     `yaml:"num_resources"`  // Number of resources
	NumTasks      int     `yaml:"num_tasks"`      // Number of tasks
	TotalUtil     float64 `yaml:"total_util"`     // Total utilization of all tasks
	MinPeriod     float64 `yaml:"min_period"`     // Period range for tasks
	MaxPeriod     float64 `yaml:"max_period"`     // Period range for tasks
	WCETRatio     float64 `yaml:"wcet_ratio"`     // Ratio of WCET2 to WCET1 for high-criticality tasks
	HighRatio     float64 `yaml:"high_ratio"`     // Probability of a task being high-criticality
	ResourceRatio float64 `yaml:"resource_ratio"` // Probability of a task using a resource
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
