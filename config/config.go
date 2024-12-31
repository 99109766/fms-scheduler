package config

import (
	"os"

	"github.com/pingcap/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	NumResources  int        `yaml:"num_resources"`  // Number of resources
	NumTasks      int        `yaml:"num_tasks"`      // Number of tasks
	TotalUtility  float64    `yaml:"total_utility"`  // Total utilization of all tasks
	PeriodRange   [2]float64 `yaml:"period_range"`   // Period range for tasks
	WCETRatio     [2]float64 `yaml:"wcet_ratio"`     // Ratio of WCET2 to WCET1 for high-criticality tasks
	HighRatio     float64    `yaml:"high_ratio"`     // Probability of a task being high-criticality
	ResourceUsage float64    `yaml:"resource_usage"` // Probability of a task using a resource
	CSFactor      float64    `yaml:"cs_factor"`      // Critical section factor
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

	if cfg.NumResources <= 0 {
		return nil, errors.New("num_resources must be greater than 0")
	}
	if cfg.NumTasks <= 0 {
		return nil, errors.New("num_tasks must be greater than 0")
	}
	if cfg.TotalUtility <= 0 || cfg.TotalUtility > 1 {
		return nil, errors.New("total_utility must be in the range (0, 1]")
	}
	if cfg.PeriodRange[0] <= 0 {
		return nil, errors.New("period_range must contain positive values")
	} else if cfg.PeriodRange[1] <= cfg.PeriodRange[0] {
		return nil, errors.New("period_range[1] must be greater than period_range[0]")
	}
	if cfg.WCETRatio[0] < 0 {
		return nil, errors.New("wcet_ratio must contain non-negative values")
	} else if cfg.WCETRatio[1] < cfg.WCETRatio[0] {
		return nil, errors.New("wcet_ratio[1] must be greater than wcet_ratio[0]")
	}
	if cfg.HighRatio < 0 || cfg.HighRatio > 1 {
		return nil, errors.New("high_ratio must be in the range [0, 1]")
	}
	if cfg.ResourceUsage < 0 || cfg.ResourceUsage > 1 {
		return nil, errors.New("resource_usage must be in the range [0, 1]")
	}
	if cfg.CSFactor <= 0 || cfg.CSFactor > 1 {
		return nil, errors.New("cs_factor must be in the range (0, 1]")
	}

	return &cfg, nil
}
