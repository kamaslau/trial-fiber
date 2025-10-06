package models

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func Seeding(DBClient *gorm.DB) {
	log.Print("Starting database seeding...")

	// 先创建权限、角色
	if err := DBClient.Create(&PermissionSeeds).Error; err != nil {
		panic("failed to seed permissions: " + err.Error())
	}
	if err := DBClient.Create(&RoleSeeds).Error; err != nil {
		panic("failed to seed roles: " + err.Error())
	}

	// 然后创建用户、账户
	if err := DBClient.Create(&UserSeeds).Error; err != nil {
		panic("failed to seed users: " + err.Error())
	}
	if err := DBClient.Create(&AccountSeeds).Error; err != nil {
		panic("failed to seed accounts: " + err.Error())
	}

	log.Print("✅ Database entity models seeded")

	// 为超级管理员分配 sa 角色
	log.Print("Assigning super admin role...")
	var superAdmin User
	var saRole Role

	if err := DBClient.Where("name = ?", "TheSuperAdminer").First(&superAdmin).Error; err != nil {
		panic("failed to find super admin: " + err.Error())
	}

	if err := DBClient.Where("name = ?", "sa").First(&saRole).Error; err != nil {
		panic("failed to find sa role: " + err.Error())
	}

	if superAdmin.ID != 0 && saRole.ID != 0 {
		if err := DBClient.Model(&superAdmin).Association("Roles").Append(&saRole); err != nil {
			panic("failed to assign role: " + err.Error())
		}
	} else {
		log.Print("Super admin or sa role not found, skipping role assignment")
	}

	// Do not have to (yet still could) connect any permissions to role, because SA explictly passes all RBAC (with proper monitoring recommended) check.
	log.Print("✅ Database relations seeded")
	log.Print("✅ All database seeding completed")
}

// [Optional] https://gorm.io/docs/migration.html#Auto-Migration
func AutoMigrate(DBClient *gorm.DB) {
	log.Print("Starting database migration...")

	// 第一步：并行迁移基础表（Permission, Role, Account）
	baseTables := []struct {
		model interface{}
		name  string
	}{
		{&Permission{}, "Permission"},
		{&Role{}, "Role"},
		{&Account{}, "Account"},
	}

	baseErrChan := make(chan error, len(baseTables))

	// 启动并行迁移
	for _, table := range baseTables {
		go func(model interface{}, name string) {
			baseErrChan <- smartMigrate(DBClient, model, name)
		}(table.model, table.name)
	}

	// 等待所有基础表迁移完成
	for i := 0; i < len(baseTables); i++ {
		if err := <-baseErrChan; err != nil {
			panic("failed to migrate database: " + err.Error())
		}
	}

	// 第二步：迁移 User 表（依赖基础表）
	if err := smartMigrate(DBClient, &User{}, "User"); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	// 第三步：并行迁移业务表（Post）
	businessTables := []struct {
		model interface{}
		name  string
	}{
		{&Post{}, "Post"},
	}

	businessErrChan := make(chan error, len(businessTables))

	// 启动并行迁移
	for _, table := range businessTables {
		go func(model interface{}, name string) {
			businessErrChan <- smartMigrate(DBClient, model, name)
		}(table.model, table.name)
	}

	// 等待所有业务表迁移完成
	for i := 0; i < len(businessTables); i++ {
		if err := <-businessErrChan; err != nil {
			panic("failed to migrate database: " + err.Error())
		}
	}

	log.Print("✅ All database migrations completed successfully!")
}

func smartMigrate(DBClient *gorm.DB, model interface{}, tableName string) error {
	log.Printf("🔄 Migrating table %s...", tableName)

	// GORM AutoMigrate 会自动：
	// 1. 检查表是否存在
	// 2. 如果不存在，创建表
	// 3. 如果存在，检查结构变化并更新
	if err := DBClient.AutoMigrate(model); err != nil {
		log.Printf("❌ Failed to migrate table %s: %v", tableName, err)
		return fmt.Errorf("failed to migrate table %s: %w", tableName, err)
	}

	log.Printf("✅ Table %s migration completed", tableName)
	return nil
}
