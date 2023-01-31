package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string  `json:"name" validate:"required,min=3,max=32"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Books    []Book  `json:"books" validate:"dive"`
	Orders   []Order `json:"orders" validate:"dive"`
}
