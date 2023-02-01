package models

import (
	"errors"

	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	BookID   uint    `json:"bookId" validate:"required"`
	CartID   uint    `json:"cartId" validate:"required"`
	Book     Book    `json:"book" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:BookID;references:ID" validate:"-"`
	Cart     Cart    `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:CartID;references:ID" validate:"-"`
	Quantity int     `json:"quantity" validate:"required,gt=0"`
	Price    float64 `json:"-" validate:"-"`
}

func (cartItem *CartItem) BeforeCreate(db *gorm.DB) error {
	book := new(Book)
	db.First(book, cartItem.BookID)
	if book.ID == 0 {
		return errors.New("invalid book id")
	}
	cartItem.Price = float64(cartItem.Quantity) * float64(book.Price)
	return nil
}
