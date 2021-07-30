package db

import (
	"fmt"
	"log"
	"os"

	"github.com/allgoodworks/Cryptobetting-Golang/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// config
var (
	DB_NAME     = ""
	DB_HOST     = ""
	DB_USER     = ""
	DB_PASSWORD = ""
	DB_PORT     = ""
	HTTP_PORT   = ""
)

func Init() {
	loadEnvirontment()
	migrateSchema()
}

func loadEnvirontment() {
	// init config
	godotenv.Load()

	// init values
	DB_NAME = os.Getenv("DB_NAME")
	DB_HOST = os.Getenv("DB_HOST")
	DB_USER = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_PORT = os.Getenv("DB_PORT")
	HTTP_PORT = os.Getenv("HTTP_PORT")
}

// schema migration
func migrateSchema() {

	db := CreateConnection()

	// db logs
	db.Debug()

	// auto migration
	db.AutoMigrate(&models.Fixture{})
	log.Println("Successfully migrated!")
}

// connection instance
func CreateConnection() *gorm.DB {

	// data source
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	// db connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}

	return db
}
