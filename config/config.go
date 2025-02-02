package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	NumResources  int        `yaml:"num_resources" validate:"gt=0"`
	NumTasks      int        `yaml:"num_tasks" validate:"gt=0"`
	TotalUtility  float64    `yaml:"total_utility" validate:"gt=0,lte=1"`
	PeriodRange   [2]float64 `yaml:"period_range" validate:"min=0,ltcsfield=PeriodRange[1]"`
	WCETRatio     [2]float64 `yaml:"wcet_ratio" validate:"min=0,ltcsfield=WCETRatio[1]"`
	HighRatio     float64    `yaml:"high_ratio" validate:"gte=0,lte=1"`
	ResourceUsage float64    `yaml:"resource_usage" validate:"gte=0,lte=1"`
	CSFactor      float64    `yaml:"cs_factor" validate:"gt=0,lte=1"`
	SimTime       float64    `yaml:"sim_time" validate:"gt=0"`
}

// LoadConfig reads the YAML configuration file from the given path and returns a pointer to a Config struct.
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
