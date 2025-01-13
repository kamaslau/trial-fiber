package drivers

import (
	"log"
	"os"
	"strconv"

	"app/src/models"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBClient 连接单例
var DBClient *gorm.DB

// ConnectDB 建立连接
func ConnectDB() {
	DB_DSN := os.Getenv("DB_DSN")
	// log.Print("DB_DSN: ", DB_DSN)

	db, err := gorm.Open(mysql.Open(DB_DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	} else {
		DBClient = db
		log.Print("👍 Database connected")
	}

	//https://gorm.io/docs/migration.html#Auto-Migration
	err = DBClient.AutoMigrate(&models.Post{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	} else {
		logOnConnected()
	}
}

// logOnConnected 写日志：数据库连接成功
func logOnConnected() {
	var payload = models.Post{
		UUID:    uuid.NewString(),
		Name:    "DB connected",
		Content: "This is an auto generated message on database connection succeed.",
	}

	result := DBClient.Create(&payload)
	if result.RowsAffected == 1 {
		log.Print("└ Succeed inserting data ID:" + strconv.Itoa(int(payload.ID)))
	} else {
		log.Print(result.Error)
		panic("failed to insert data")
	}
}
