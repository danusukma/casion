package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	TransactionID string         `json:"transaction_id" gorm:"primaryKey;type:char(36)"`
	UserID        string         `json:"user_id" gorm:"not null"`
	Type          string         `json:"type" gorm:"not null"`
	Status        string         `json:"status" gorm:"not null"`
	Amount        float64        `json:"amount" gorm:"not null"`
	Remarks       string         `json:"remarks"`
	BalanceBefore float64        `json:"balance_before" gorm:"not null"`
	BalanceAfter  float64        `json:"balance_after" gorm:"not null"`
	TargetUserID  string         `json:"target_user_id"`
	CreatedAt     time.Time      `json:"created_date"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
