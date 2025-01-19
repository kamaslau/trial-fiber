package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UUID string `json:"uuid" gorm:"index;type:varchar(36)"`
	Name string `json:"name"`

	Content string  `json:"content"`
	Excerpt *string `json:"excerpt" gorm:"type:varchar(255)"` // Optional field name, * means nullable
}
