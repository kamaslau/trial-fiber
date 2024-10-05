package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UUID string `gorm:"index;type:varchar(36)"`
	Name string

	Content string
	Excerpt *string `gorm:"type:varchar(255)"` // Optional field name
}
