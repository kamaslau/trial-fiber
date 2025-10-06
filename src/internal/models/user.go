package models

import (
	"app/src/internal/utils/uuid"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID string `json:"uuid" gorm:"<-:create;unique;uniqueIndex;type:varchar(36)"`
	Name string `json:"name" gorm:"unique;uniqueIndex;type:varchar(255)"`

	NickName *string `json:"nickname" gorm:"unique;uniqueIndex;type:varchar(255)"`
	Avatar   *string `json:"avatar,omitempty" gorm:"type:varchar(255)"`
	Password *string `json:"password,omitempty" gorm:"type:char(64)"`

	Roles []Role `json:"-" gorm:"many2many:user_roles;"`
}

type UserFieldsCreate struct {
	Name        string `json:"name" validate:"required"`
	NickName    string `json:"nickname"`
	Description string `json:"description"`
}

type UserFieldsUpdate struct {
	NickName    string `json:"nickname"`
	Description string `json:"description"`
}

var UserSeeds = []User{
	{
		UUID: uuid.NewString(),
		Name: "TheSuperAdminer",
	},
	{
		UUID:     uuid.NewString(),
		Name:     "FirstUser",
		Password: stringPtr("$2a$12$GAJZK4iROhR5iGiJigF0I.Rju5YFLOmvRGfnpDM.QkB7izl7mx2o2"), // 12 cost bcrypted 123abc456def
	},
}
