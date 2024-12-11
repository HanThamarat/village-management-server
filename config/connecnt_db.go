package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"tectnexify.github.com/e-payment/models"
)

func DatabaseConnection() *gorm.DB {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading database .env file")
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))

	var (
		host     = os.Getenv("DB_HOST")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
	)

	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{})

	db.AutoMigrate(
		&models.Users{},
		&models.Roles{},
		&models.BankCredentials{},
		&models.VallageHouses{},
		&models.VallageOwnerShips{},
	)

	if err != nil {
		panic(err)
	}

	return db
}
