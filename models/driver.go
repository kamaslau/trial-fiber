package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

// DBClient 连接单例
var DBClient *gorm.DB

// Connect 建立连接
func Connect() {
	dsn := os.Getenv("DSN")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DBClient = db
	fmt.Print("DB connected\n")
}
