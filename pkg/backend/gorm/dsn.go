package gorm

import (
	"errors"
	"regexp"
)

type DSNConfig struct {
	Driver    string
	Connexion string
}

var DSNRegex = regexp.MustCompile(`(sqlite|mysql|postgres|sqlserver)://(.*)`)
var ErrDSNParseError = errors.New("unable to parse dsn")

func ParseDSN(dsn string) (DSNConfig, error) {
	r := DSNRegex.FindStringSubmatch(dsn)
	cfg := DSNConfig{}
	if r == nil {
		return cfg, ErrDSNParseError
	}
	cfg.Driver = r[1]
	cfg.Connexion = r[2]
	return cfg, nil
}
