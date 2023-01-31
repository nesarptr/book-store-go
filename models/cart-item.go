package models

import (
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	BookID   uint `json:"bookId"`
	CartID   uint `json:"cartId"`
	Book     Book `json:"book" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Cart     Cart `json:"cart" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Quantity int  `json:"quantity"`
}
