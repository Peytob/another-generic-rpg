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
	Log    *LogConfiguration    `env-prefix:"LOG_"    validate:"omitempty" yaml:"log"`
	Window *WindowConfiguration `env-prefix:"WINDOW_" validate:"omitempty" yaml:"window"`
}

type ServerConfiguration struct {
	Log  *LogConfiguration  `env-prefix:"LOG_"  validate:"omitempty" yaml:"log"`
	HTTP *HTTPConfiguration `env-prefix:"HTTP_" validate:"omitempty" yaml:"http"`
}

type LogConfiguration struct {
	Enabled bool       `env:"ENABLED" yaml:"enabled"`
	Level   slog.Level `env:"LEVEL"   yaml:"level"`
}

type HTTPConfiguration struct {
	Port int    `env:"PORT" validate:"required,min=1,max=65535" yaml:"port"`
	Host string `env:"HOST" validate:"required"                 yaml:"host"`
}

type WindowConfiguration struct {
	Width  int `env:"WIDTH"  validate:"required,min=320,max=7680" yaml:"width"`
	Height int `env:"HEIGHT" validate:"required,min=240,max=4800" yaml:"height"`
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
