package config

import (
	"os"
	"fmt"
	"time"
	toml "github.com/pelletier/go-toml/v2"
)

type AuthConfig struct {
	AuthKeyTTL 	int64			 `toml:"auth_key_ttl"`
	CleanupDelay string			 `toml:"cleanup_delay"`
	CleanupDuration time.Duration `toml:"-"`
}

type ServerConfig struct {
	Host       string            `toml:"host"` // Host to listen to
	LimiterMax int				 `toml:"limiter_max"` // Max request during window time for rate limiter in seconds
	LimiterWindow int 			 `toml:"limiter_window"` // Window time for rate limiter in seconds
	LoginLimiterMax int			 `toml:"login_limiter_max"` // Max request during window time for rate limiter in seconds
	LoginLimiterWindow int		 `toml:"login_limiter_window"` // Window time for rate limiter in seconds
}

type DBConfig struct {
	DSN        string            `toml:"dsn"`  // Database DSN sqlite://file, po
	Debug	   bool				 `toml:"debug"` // Debug Mode for Grom Db
}

type AppConfig struct {
	Server     	ServerConfig 	 `toml:"server"`
	DB		   	DBConfig 		 `toml:"db"`
	Auth		AuthConfig       `toml:"auth"`
	SurveyPath 	string            `toml:"survey_path"`
	Users      map[string]string `toml:"users"` // User : Argon Hash
}

func parseDuration(value string, def time.Duration) (time.Duration, error) {
	if(value == "") {
		return def, nil
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		return time.Duration(0), fmt.Errorf("invalid time duration '%s' : %s", value, err.Error())
	}
	if(d <= time.Duration(0)) {
		return def, nil
	}
	return d, nil
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

	cleanupDuration, err := parseDuration(cfg.Auth.CleanupDelay, 24 * time.Hour)

	if(cfg.Auth.AuthKeyTTL <= 0) {
		cfg.Auth.AuthKeyTTL = 60 * 60
	}

	cfg.Auth.CleanupDuration = cleanupDuration

	return &cfg, nil
}

func (self *AppConfig) Show() {
	fmt.Printf("Survey path=%s\n", self.SurveyPath)
	fmt.Printf("Using authKey TTL=%d Cleanup Delay=%s\n", self.Auth.AuthKeyTTL, self.Auth.CleanupDuration)
	if(self.DB.Debug) {
		fmt.Printf("Database debug mode is On\n", self.SurveyPath)
	}
}