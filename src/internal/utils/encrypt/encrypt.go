package encrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// hash returns the bcrypt string (cost 12) for a plain password.
func BcryptHash(plain string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(plain), 12)
	return string(hashBytes), err
}

// ok reports whether the plain password matches the stored bcrypt string.
func BcryptValidate(plain, storedHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(plain)) == nil
}
