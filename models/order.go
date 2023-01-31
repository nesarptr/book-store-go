package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	TotalPrice float64     `json:"totalPrice"`
	Books      []OrderItem `json:"books"`
	UserID     uint        `json:"userId"`
	Owner      User        `json:"owner" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type OrderItem struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}
