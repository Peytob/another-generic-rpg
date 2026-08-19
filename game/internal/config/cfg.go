package config

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

type ClientConfiguration struct {
	Log    *LogConfiguration    `yaml:"log" env-prefix:"LOG_" validate:"omitempty"`
	Window *WindowConfiguration `yaml:"window" env-prefix:"WINDOW_" validate:"omitempty"`
}

type ServerConfiguration struct {
	Log  *LogConfiguration  `yaml:"log" env-prefix:"LOG_" validate:"omitempty"`
	HTTP *HTTPConfiguration `yaml:"http" env-prefix:"HTTP_" validate:"omitempty"`
}

type LogConfiguration struct {
	Enabled bool       `yaml:"enabled" env:"ENABLED"`
	Level   slog.Level `yaml:"level" env:"LEVEL"`
}

type HTTPConfiguration struct {
	Port int    `yaml:"port" env:"PORT" validate:"required,min=1,max=65535"`
	Host string `yaml:"host" env:"HOST" validate:"required"`
}

type WindowConfiguration struct {
	Width  int `yaml:"width" env:"WIDTH" validate:"required,min=320,max=7680"`
	Height int `yaml:"height" env:"HEIGHT" validate:"required,min=240,max=4800"`
}

func LoadConfiguration[T ClientConfiguration | ServerConfiguration](path string) (*T, error) {
	cfg := new(T)

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, fmt.Errorf("failed to load configuration %s from and envs: %w", path, err)
	}

	if err := validateConfiguration(cfg); err != nil {
		return nil, fmt.Errorf("failed client configuration validation: %w", err)
	}

	return cfg, nil
}

func validateConfiguration[T ClientConfiguration | ServerConfiguration](cfg *T) error {
	validate := validator.New()

	err := validate.Struct(cfg)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			fields := make([]string, 0, len(validateErrs))
			for _, e := range validateErrs {
				fields = append(fields, e.Namespace())
			}

			return fmt.Errorf("configuration validation failed: %w: %s", err, strings.Join(fields, "; "))
		}

		return fmt.Errorf("configuration validation failed: %w", err)
	}

	return nil
}
