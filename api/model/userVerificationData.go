package model

import (
	"time"

	"gorm.io/gorm"
)

type VerificationData struct {
	gorm.Model
    Email string `json:"email" gorm:"not null"`
    VerificationToken string `json:"token" gorm:"not null"`
    ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
}