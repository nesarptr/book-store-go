package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nesarptr/book-store-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() error {
	HOST, err := GetEnv("HOST")
	USER, _ := GetEnv("USER")
	PASSWORD, _ := GetEnv("PASSWORD")
	DATABASE, _ := GetEnv("DATABASE")
	PORT, _ := GetEnv("PORT")
	SSLMODE, _ := GetEnv("SSLMODE")
	if err != nil {
		return err
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", HOST, USER, PASSWORD, DATABASE, PORT, SSLMODE)
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	d.AutoMigrate(&models.User{}, &models.Book{}, &models.Order{}, &models.Cart{}, &models.CartItem{})
	db = d
	return nil
}

func GetDB() *gorm.DB {
	return db
}

func GetEnv(key string) (string, error) {
	err := godotenv.Load("E:\\Projects\\Portfolio-Projects\\book-store-go\\.env")
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return os.Getenv(key), nil
}
