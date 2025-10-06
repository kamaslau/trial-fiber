package models

import (
	"app/src/internal/utils/uuid"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	UUID string `json:"uuid" gorm:"<-:create;unique;uniqueIndex;type:varchar(36)"`

	UserID   uint   `json:"userId" gorm:"<-:create;index;uniqueIndex:user_provider_idn;type:uint"`
	Provider string `json:"provider" gorm:"<-:create;index;uniqueIndex:user_provider_idn;type:varchar(100)"`
	Idn      string `json:"idn" gorm:"<-:create;index;uniqueIndex:user_provider_idn;type:varchar(255);comment:Identification Code Or Account Number"` // IDN, abbr. for Identification Number.
}

type AccountFieldsCreate struct {
	// get user_id from token
	UserID   uint   `json:"userId" validate:"required"`
	Provider string `json:"provider" validate:"required"`
	Idn      string `json:"idn" validate:"required"`
}

func AccountValidateCreate(input AccountFieldsCreate) (bool, error) {
	validate := validator.New()

	if err := validate.Struct(input); err != nil {
		// Validation failed, handle the error
		errors := err.(validator.ValidationErrors)
		return false, errors
	}

	return true, nil
}

type AccountFieldsUpdate struct {
	// Do not allow update
}

var AccountSeeds = []Account{
	{
		UUID:     uuid.NewString(),
		UserID:   1,
		Provider: "mobile",
		Idn:      "19988889999",
	},
	{
		UUID:     uuid.NewString(),
		UserID:   2,
		Provider: "email",
		Idn:      "lucky@me.com",
	},
}
