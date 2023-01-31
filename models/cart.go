package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	TotalPrice float64    `json:"totalPrice" validate:"required,gt=0"`
	Books      []CartItem `json:"books" validate:"dive"`
	UserID     uint       `json:"userId" validate:"required"`
	Owner      User       `json:"owner" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID" validate:"dive"`
}
