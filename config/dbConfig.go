package config

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var myDb *gorm.DB

func Connect() {
	LoadEnv()

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Error while connecting to the database: ", err)
	}
	myDb = db
}

func GetDb() *gorm.DB {
	return myDb
}
