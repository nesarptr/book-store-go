package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	TotalPrice float64     `json:"totalPrice" validate:"required,gt=0"`
	Books      []OrderItem `json:"books" validate:"dive"`
	UserID     uint        `json:"userId" validate:"required"`
	Owner      User        `json:"owner" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID" validate:"dive"`
}

type OrderItem struct {
	gorm.Model
	OrderID uint   `json:"orderId" validate:"required"`
	Order   Order  `json:"order" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:OrderID;references:ID" validate:"dive"`
	Title   string `json:"title" validate:"required,min=3,max=32"`
}
