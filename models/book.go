package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title       string  `json:"title"`
	Price       float32 `json:"price"`
	ImgUrl      string  `json:"imgUrl"`
	Description string  `json:"description"`
	UserID      uint    `json:"userId"`
	Author      User    `json:"author" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
