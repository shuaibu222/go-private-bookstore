package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var myDb *gorm.DB

func Connect() {

	dsn := "shuaibu:Shuaibu%12345%@tcp(localhost:3306)/online_books?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Error while connecting to the database: ", err)
	}
	myDb = db
}

func GetDb() *gorm.DB {
	return myDb
}
