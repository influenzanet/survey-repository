package utils

import (
	"github.com/influenzanet/user-management-service/pkg/pwhash"
)

func CheckPassword(hash string, password string) (bool, error) {
	return pwhash.ComparePasswordWithHash(hash, password)
}

func HashPassword(password string) (string, error) {
	return pwhash.HashPassword(password)
}
