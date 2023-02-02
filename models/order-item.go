package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	OrderID     uint    `json:"orderId" validate:"required"`
	Order       Order   `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:OrderID;references:ID" validate:"-"`
	Title       string  `json:"title" validate:"required,min=3,max=32"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	ImgUrl      string  `json:"imgUrl" validate:"required,min=10"`
	Description string  `json:"description"`
	BookID      uint    `json:"bookId" validate:"required"`
	UserID      uint    `json:"owner" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,gt=0"`
}
