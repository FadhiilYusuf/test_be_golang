package database

import (
	"fmt"
	"log"
	"os"

	"github.com/cngJo/golang-api-auth/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Instance *gorm.DB
var dbError error

// Fungsi untuk membangun connection string berdasarkan environment variables
func buildConnectionString() string {
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)
}

func Connect() {
	// Gunakan connection string yang di-build dari environment variables
	connectionString := buildConnectionString()

	Instance, dbError = gorm.Open(mysql.Open(connectionString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if dbError != nil {
		log.Fatal("Failed to connect to database: ", dbError)
		panic("Database connection failed")
	}

	log.Println("Database connected")
}

func Migrate() {
	Instance.AutoMigrate(&models.User{})
	Instance.AutoMigrate(&models.Product{})
	Instance.AutoMigrate(&models.Order{})
	log.Println("Database migrated")
}
