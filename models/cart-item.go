package models

import (
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	BookID   uint `json:"bookId" validate:"required"`
	CartID   uint `json:"cartId" validate:"required"`
	Book     Book `json:"book" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:BookID;references:ID" validate:"-"`
	Cart     Cart `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:CartID;references:ID" validate:"-"`
	Quantity int  `json:"quantity" validate:"required,gt=0"`
}
