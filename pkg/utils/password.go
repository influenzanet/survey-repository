package utils

import (
	"github.com/alexedwards/argon2id"
)

func CheckPassword(hash string, password string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}
