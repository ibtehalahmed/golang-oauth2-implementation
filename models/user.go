package models

import (
	"time"
)

type User struct {
	ID         uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name   	string    `gorm:"not null"json:"name" form:"name"`
	Surname   string    `json:"surname" form:"surname"`
	Username   string    `gorm:"unique; index" json:"username" form:"username"`

	Email      string    `gorm:"not null;unique" json:"email" form:"email"`
	Phone      string    `json:"phone" form:"phone"`
	Active      bool    `json:"active" form:"active"`

	Password   string    `gorm:"size:500;not null;" json:"password" form:"password"`
	PasswordConfirmation   string    `gorm:"-" json:"password_confirmation" form:"password_confirmation" db:"-"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	EmailVerifiedAt  time.Time `json:"email_verified_at"`
}