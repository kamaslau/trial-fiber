package auth

import (
	"errors"
	"log"

	"app/src/internal/models"
	"app/src/internal/utils/drivers"
	"app/src/internal/utils/encrypt"
)

type AuthInput struct {
	Provider string `json:"provider"`
	Idn      string `json:"idn"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

// Auth with password
func AuthByPassword(userID uint, password string) error {
	// log.Printf("AuthByPassword: userID=%d, password=%s", userID, password) // test ONLY

	var err error

	// Fetch User by id
	var filter = map[string]any{
		"id": userID,
	}
	var user models.User
	if err := drivers.DBClient.Where(filter).First(&user).Error; err != nil {
		log.Printf("FindOne: database query failed: %v", err)
		return err
	}

	// Check password
	if user.Password == nil {
		log.Printf("User password is nil")
		return errors.New("user password not set")
	}

	// Validate password
	// log.Printf("Validating password: plain='%s', hash='%s'", password, *user.Password) // test ONLY
	if isValid := encrypt.BcryptValidate(password, *user.Password); !isValid {
		log.Printf("Password validation failed")
		err = errors.New("password match failed")
	}

	return err
}
