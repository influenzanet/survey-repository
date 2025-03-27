package config

import (
	"os"

	toml "github.com/pelletier/go-toml/v2"
)

type ServerConfig struct {
	Host       string            `toml:"host"` // Host to listen to
	LimiterMax int				 `toml:"limiter_max"` // Max request during window time for rate limiter in seconds
	LimiterWindow int 			 `toml:"limiter_window"` // Window time for rate limiter in seconds
}

type AppConfig struct {
	Server     ServerConfig `toml:"server"`
	DSN        string            `toml:"dsn"`  // Database DSN sqlite://file, po
	SurveyPath string            `toml:"survey_path"`
	Users      map[string]string `toml:"users"` // User : Argon Hash
}

func getEnvOr(name string, def string) string {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	return v
}

func LoadConfig() (*AppConfig, error) {

	file := getEnvOr("APP_CONFIG", "app.toml")

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var cfg AppConfig

	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
