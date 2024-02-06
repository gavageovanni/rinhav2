package config

import (
	"baitadores-rinhav2/transaction"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB
var e error

func DatabaseInit() {
	host := "localhost"
	user := "rinhav2"
	password := "rinhav2"
	dbName := "rinhav2"
	port := 5432

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, dbName, port)
	database, e = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if e != nil {
		panic(e)
	}
	database.AutoMigrate(&transaction.TransactionModel{}, &transaction.UserModel{})

}

func DB() *gorm.DB {
	return database
}
