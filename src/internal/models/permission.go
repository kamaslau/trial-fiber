package models

import (
	"app/src/internal/utils/uuid"

	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	UUID string `json:"uuid" gorm:"<-:create;unique;uniqueIndex;type:varchar(36)"`
	Name string `json:"name" gorm:"unique;uniqueIndex;type:varchar(100)"`

	Description *string `json:"description,omitempty" gorm:"type:varchar(255)"`
}

type PermissionFieldsCreate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type PermissionFieldsUpdate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

var PermissionSeeds = []Permission{
	// User permissions
	{
		UUID:        uuid.NewString(),
		Name:        "user::read::all",
		Description: stringPtr("User->Read->All"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "user::read::self",
		Description: stringPtr("User->Read->Self"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "user::write::all",
		Description: stringPtr("User->Read->All"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "user::create::self",
		Description: stringPtr("User->Create->Self"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "user::update::self",
		Description: stringPtr("User->Update->Self"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "user::delete::self",
		Description: stringPtr("User->Delete->Self"),
	},

	// Role permissions
	{
		UUID:        uuid.NewString(),
		Name:        "role::read::all",
		Description: stringPtr("Role->Read->All"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "role::read::self",
		Description: stringPtr("Role->Read->Self"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "role::write::all",
		Description: stringPtr("User->Read->All"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "role::create::self",
		Description: stringPtr("Role->Create->Self"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "role::update::self",
		Description: stringPtr("Role->Update->Self"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "role::delete::self",
		Description: stringPtr("Role->Delete->Self"),
	},

	// Permission permissions
	{
		UUID:        uuid.NewString(),
		Name:        "permission::read::all",
		Description: stringPtr("Permission->Read->All"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "permission::read::self",
		Description: stringPtr("Permission->Read->Self"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "permission::write::all",
		Description: stringPtr("User->Read->All"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "permission::create::self",
		Description: stringPtr("Permission->Create->Self"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "permission::update::self",
		Description: stringPtr("Permission->Update->Self"),
	},
	{
		UUID:        uuid.NewString(),
		Name:        "permission::delete::self",
		Description: stringPtr("Permission->Delete->Self"),
	},
}
