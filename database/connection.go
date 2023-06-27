package database

import (
	"fmt"
	"log"
	"os"

	"github.com/adityarudrawar/go-backend/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	
	dsn := get_dsn()
	
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("connection to database failed")
	}
	
	// Auto create tables with models by using the AutoMigrate method 
	connection.AutoMigrate(&models.User{})

	DB = connection
}

func get_dsn() string {

	err := godotenv.Load()
	if err != nil {
		log.Panic("Failed to load ENV")
	}
	POSTGRES_HOST := os.Getenv("POSTGRES_HOST")
	POSTGRES_PORT := os.Getenv("POSTGRES_PORT")

	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DBNAME := os.Getenv("POSTGRES_DBNAME")

	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DBNAME)
}