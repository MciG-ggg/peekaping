package utils

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateConfig unmarshals the config string into the given struct and validates it using struct tags.
func ValidateConfig[T any](config string) (*T, error) {
	var cfg T
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return nil, err
	}
	if err := validate.Struct(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
