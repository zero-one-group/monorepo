package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("error hashing: %v\n", err)
		return "", err
	}
	hash := string(hashedBytes)
	return hash, nil
}

func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("Password is invalid: %v\n", password)
		return false
	}
	 return true
}
