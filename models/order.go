package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	TotalPrice float64     `json:"totalPrice" validate:"required,gt=0"`
	Books      []OrderItem `json:"books" validate:"dive"`
	UserID     uint        `json:"-" validate:"required"`
	Owner      User        `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID" validate:"-"`
}
