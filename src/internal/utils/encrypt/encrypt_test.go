package encrypt

import (
	"testing"
)

func TestBcryptHashAndValidate(t *testing.T) {
	plain := "123abc456def"

	hash, err := BcryptHash(plain)
	if err != nil {
		t.Fatalf("Failed to Hash Bcrypt: %v", err)
	}

	// Test correct password
	if BcryptValidate(plain, hash) != true {
		t.Fatal("Failed to Validate Bcrypt")
	}

	// Test wrong password
	hash = "some_made_up_string"
	if BcryptValidate(plain, hash) == true {
		t.Fatal("Failed to Validate Bcrypt")
	}
}
