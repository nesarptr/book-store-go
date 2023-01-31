package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	TotalPrice float64    `json:"totalPrice"`
	Books      []CartItem `json:"books"`
	UserID     uint       `json:"userId"`
	Owner      User       `json:"owner"`
}
