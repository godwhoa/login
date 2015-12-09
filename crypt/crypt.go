package crypt

import (
	"golang.org/x/crypto/bcrypt"
)

func Check(hashed string, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
	if err == nil {
		return true
	} else {
		return false
	}
}

func Hash(pass string) string {
	password := []byte(pass)

	// Hashing the password with the default cost of 10
	hashed, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashed)
}
