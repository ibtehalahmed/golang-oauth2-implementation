package models

import (
	"time"
)

type Oauth2Client struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Secret     string    `json:"secret"`
	UserID       uint64         `json:"user_id"`
	Name   string         `json:"name"`
	Redirect string         `json:"redirect"`
	Domain string         `json:"domain"`
	Data string         `json:"data"`
	PersonalAccessClient string         `json:"personal_access_client"`
	PasswordClient string         `json:"password_client"`
	Revoked  string `json:"revoked"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}