package config

import (
	"os"
	"reflect"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	NumResources  int        `yaml:"num_resources" validate:"min=0"`
	NumTasks      int        `yaml:"num_tasks" validate:"min=0"`
	TotalUtility  float64    `yaml:"total_utility" validate:"min=0,max=1"`
	PeriodRange   [2]float64 `yaml:"period_range" validate:"min=0,valid_range"`
	DeadlineRatio [2]float64 `yaml:"deadline_ratio" validate:"min=0,valid_range"`
	WCETRatio     [2]float64 `yaml:"wcet_ratio" validate:"min=0,valid_range"`
	HighRatio     float64    `yaml:"high_ratio" validate:"min=0,max=1"`
	ResourceUsage [2]int     `yaml:"resource_usage" validate:"min=0,ltefield=NumResources,valid_range"`
	CSFactor      float64    `yaml:"cs_factor" validate:"min=0,max=1"`
	CSRange       [2]int     `yaml:"cs_range" validate:"min=0,valid_range"`
	SimulateTime  float64    `yaml:"simulation_time" validate:"min=0"`
}

func defineValidators(validate *validator.Validate) {
	// valid_range checks that the first element is ≤ the second element.
	validate.RegisterValidation("valid_range", func(fl validator.FieldLevel) bool {
		if fl.Field().Kind() != reflect.Array || fl.Field().Len() != 2 {
			return false
		}

		switch fl.Field().Index(0).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fl.Field().Index(0).Int() <= fl.Field().Index(1).Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return fl.Field().Index(0).Uint() <= fl.Field().Index(1).Uint()
		case reflect.Float32, reflect.Float64:
			return fl.Field().Index(0).Float() <= fl.Field().Index(1).Float()
		default:
			return false
		}
	})

	// ltefield checks that the field is less than or equal to the field specified by the parameter.
	validate.RegisterValidation("ltefield", func(fl validator.FieldLevel) bool {
		otherField := fl.Parent().FieldByName(fl.Param())
		switch fl.Field().Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < fl.Field().Len(); i++ {
				switch fl.Field().Index(i).Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if fl.Field().Index(i).Int() > otherField.Int() {
						return false
					}
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					if fl.Field().Index(i).Uint() > otherField.Uint() {
						return false
					}
				case reflect.Float32, reflect.Float64:
					if fl.Field().Index(i).Float() > otherField.Float() {
						return false
					}
				default:
					return false
				}
			}
			return true
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fl.Field().Int() <= otherField.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return fl.Field().Uint() <= otherField.Uint()
		case reflect.Float32, reflect.Float64:
			return fl.Field().Float() <= otherField.Float()
		default:
			return false
		}
	})
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
	defineValidators(validate)
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
