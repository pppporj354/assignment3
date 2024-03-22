package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "root:123@tcp(127.0.0.1:55001)/project?charset=utf8mb4&parseTime=True&loc=Local"
	var err error

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}
	fmt.Println("Connected to the database")
}
