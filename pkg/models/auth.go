package models

type AuthKey struct {
	Key string `json:"key" gorm:"index"`
	Created int64 `json:"created_at"`
	User string `json:"user" gorm:"index"`
}