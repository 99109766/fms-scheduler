package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	NumResources  int        `yaml:"num_resources" validate:"min=0"`
	NumTasks      int        `yaml:"num_tasks" validate:"min=0"`
	TotalUtility  float64    `yaml:"total_utility" validate:"min=0,max=1"`
	PeriodRange   [2]float64 `yaml:"period_range" validate:"min=0,valid_range"`
	WCETRatio     [2]float64 `yaml:"wcet_ratio" validate:"min=0,valid_range"`
	HighRatio     float64    `yaml:"high_ratio" validate:"min=0,max=1"`
	ResourceUsage float64    `yaml:"resource_usage" validate:"min=0,max=1"`
	CSFactor      float64    `yaml:"cs_factor" validate:"min=0,max=1"`
	SimulateTime  float64    `yaml:"sim_time" validate:"min=0"`
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
	// valid_range checks that the first element is ≤ the second element.
	validate.RegisterValidation("valid_range", func(fl validator.FieldLevel) bool {
		arr, ok := fl.Field().Interface().([2]float64)
		if !ok {
			return false
		}
		return arr[0] <= arr[1]
	})
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
