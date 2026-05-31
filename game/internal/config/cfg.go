package config

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

type ClientConfiguration struct {
	Log    *LogConfiguration    `yaml:"log" env-prefix:"LOG_" validate:"omitempty"`
	Window *WindowConfiguration `yaml:"window" env-prefix:"WINDOW_" validate:"omitempty"`
}

type ServerConfiguration struct {
	Log  *LogConfiguration  `yaml:"log" env-prefix:"LOG_" validate:"omitempty"`
	Http *HttpConfiguration `yaml:"http" env-prefix:"HTTP_" validate:"omitempty"`
}

type LogConfiguration struct {
	Enabled bool       `yaml:"enabled" env:"ENABLED"`
	Level   slog.Level `yaml:"level" env:"LEVEL"`
}

type HttpConfiguration struct {
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
			// TODO Return list of invalid fields
		}

		return err
	}

	return nil
}
