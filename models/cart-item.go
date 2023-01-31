package models

import (
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	BookID   uint `json:"bookId" validate:"required"`
	CartID   uint `json:"cartId" validate:"required"`
	Book     Book `json:"book" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" validate:"dive"`
	Cart     Cart `json:"cart" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" validate:"dive"`
	Quantity int  `json:"quantity" validate:"required,gt=0"`
}
