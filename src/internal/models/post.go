package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UUID string `json:"uuid" gorm:"<-:create;unique;uniqueIndex;type:varchar(36)"`
	Name string `json:"name"`

	Content string  `json:"content"`
	Excerpt *string `json:"excerpt,omitempty" gorm:"type:varchar(255)"`
}

type PostFieldsCreate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type PostFieldsUpdate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}
