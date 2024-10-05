package models

import (
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBClient 连接单例
var DBClient *gorm.DB

// ConnectDB 建立连接
func ConnectDB() {
	dsn := os.Getenv("DSN")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	fmt.Println("DB connected")
	DBClient = db

	logOnConnected()

	//https://gorm.io/docs/migration.html#Auto-Migration
	err = DBClient.AutoMigrate(&Post{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}

}

// logOnConnected 写日志：数据库连接成功时
func logOnConnected() {
	var payload = Post{
		UUID:    uuid.NewString(),
		Name:    "DB connected",
		Content: "This is an auto generated message on database connection succeed.",
	}
	result := DBClient.Create(&payload)
	if result.RowsAffected == 1 {
		fmt.Println("succeed to insert data: id:" + strconv.Itoa(int(payload.ID)))
	} else {
		fmt.Println(result.Error)
		panic("failed to insert data")
	}
}
