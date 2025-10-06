package models

import (
	"app/src/internal/utils/uuid"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	UUID string `json:"uuid" gorm:"<-:create;unique;uniqueIndex;type:varchar(36)"`
	Name string `json:"name" gorm:"unique;uniqueIndex;type:varchar(100)"`

	Description *string `json:"description,omitempty" gorm:"type:varchar(255)"`

	Permissions []Permission `json:"-" gorm:"many2many:role_permissions;"`
}

type RoleFieldsCreate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type RoleFieldsUpdate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

var RoleSeeds = []Role{
	{
		UUID:        uuid.NewString(),
		Name:        "sa",
		Description: stringPtr("Super Administrator"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "adminer",
		Description: stringPtr("Administrator"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "user",
		Description: stringPtr("Logged-in user"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "prisoner",
		Description: stringPtr("Blacklisted user"),
	},
}
