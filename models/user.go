package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Books    []Book  `json:"books"`
	Orders   []Order `json:"orders"`
}
