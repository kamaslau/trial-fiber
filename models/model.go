package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UUID    string // Remember to generate and store as binary(16) to save space
	Title   string
	Content string
	Excerpt *string // Optional field name
}
