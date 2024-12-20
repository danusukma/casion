package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID      string         `json:"user_id" gorm:"type:char(36);uniqueIndex"`
	FirstName   string         `json:"first_name" gorm:"not null"`
	LastName    string         `json:"last_name" gorm:"not null"`
	PhoneNumber string         `json:"phone_number" gorm:"primaryKey;type:varchar(15)"`
	Address     string         `json:"address" gorm:"not null"`
	Pin         string         `json:"-" gorm:"not null"`
	Balance     float64        `json:"balance" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_date"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
