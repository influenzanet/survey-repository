package config

import (
	"os"

	toml "github.com/pelletier/go-toml/v2"
)

type AppConfig struct {
	Host       string            `toml:"host"`
	DSN        string            `toml:"dsn"`
	SurveyPath string            `toml:"survey_path"`
	Users      map[string]string `toml:"users"`
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
