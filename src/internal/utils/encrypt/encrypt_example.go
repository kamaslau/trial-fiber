package encrypt

import (
	"fmt"
	"log"
)

func ExampleHashAndValidate() {
	pwd := "123456"

	// 1. Create (store this in DB)
	h, err := BcryptHash(pwd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("stored:", h)

	// 2. Validate
	fmt.Println("match:", BcryptValidate("123456", h)) // true
	fmt.Println("match:", BcryptValidate("wrong", h))  // false
}
