package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Name  string `json:"name" gorm:"not null"`
    Email string `json:"email" gorm:"not null; unique"`
    Password string `json:"password" gorm:"null"`
	IsSubscribed bool `json:"is_subscribed" gorm:"default:false"`
    IsVerified bool `json:"is_verified" gorm:"default:false"`
    CreatedAt  time.Time `json:"createdat" `
	UpdatedAt  time.Time `json:"updatedat"`

}
