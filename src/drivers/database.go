package drivers

import (
	"fmt"
	"log"
	"os"
	"time"

	"app/src/models"

	"github.com/google/uuid"
	// "gorm.io/driver/mysql" // FYI
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBClient 连接单例
var DBClient *gorm.DB

// ConnectDB 建立连接
func ConnectDB() {
	DB_DSN := os.Getenv("DB_DSN")
	// log.Print("DB_DSN: ", DB_DSN)

	if DB_DSN == "" {
		panic("❌ Database configs not found")
	}

	// Confs
	configs := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	}

	if db, err := gorm.Open(postgres.Open(DB_DSN), configs); err != nil {
		panic("❌ failed to connect database: " + err.Error())
	} else {
		DBClient = db
		log.Print("👍 Database connected")
	}

	// [Optional] https://gorm.io/docs/migration.html#Auto-Migration
	if err := DBClient.AutoMigrate(&models.Post{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	// logOnDBConnected()
}

// logOnDBConnected 写日志：数据库连接成功
func logOnDBConnected() {
	var payload = models.Post{
		UUID:    uuid.NewString(),
		Name:    "DB connected",
		Content: fmt.Sprintf("This is an auto generated message on database connection succeed at %s", time.Now().Format(time.RFC3339)),
	}

	if err := DBClient.Create(&payload).Error; err != nil {
		log.Printf("└ Failed to insert data: %v", err)
		return
	}
	log.Printf("└ Succeed to insert data ID: %d", payload.ID)
}

// CloseDB Terminate present database connection (if any)
func CloseDB() {
	if DBClient != nil {
		if db, err := DBClient.DB(); err == nil {
			if err := db.Close(); err != nil {
				log.Printf("❌ Failed to close database connection: %v", err)
				return
			}

			log.Print("✅ Database connection closed safely")
		}
	}
}
