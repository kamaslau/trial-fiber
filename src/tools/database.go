package main

import (
	"app/src/internal/models"
	"app/src/internal/utils"
	"app/src/internal/utils/drivers"
)

func main() {
	utils.LoadEnv()
	drivers.ConnectDB()

	// 执行数据库迁移和种子数据
	models.AutoMigrate(drivers.DBClient)
	models.Seeding(drivers.DBClient)
}
