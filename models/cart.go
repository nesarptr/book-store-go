package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	TotalPrice float64    `json:"totalPrice" validate:"required,gte=0"`
	Books      []CartItem `json:"books" validate:"-"`
	UserID     uint       `json:"-" validate:"required"`
	Owner      User       `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID" validate:"-"`
}

func (cart *Cart) Create(db *gorm.DB) error {
	if result := db.Create(cart); result.Error != nil {
		return result.Error
	}
	return nil
}
