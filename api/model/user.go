package model

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Name  string `json:"name" gorm:"not null"`
    Email string `json:"email" gorm:"not null; unique"`
    Password string `json:"password" gorm:"null"`
	IsSubscribed bool `json:"is_subscribed" gorm:"default:false"`
}