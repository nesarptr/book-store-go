package models

import (
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	BookID   uint `json:"bookId"`
	CartID   uint `json:"cartId"`
	Book     Book `json:"book"`
	Quantity int  `json:"quantity"`
}
