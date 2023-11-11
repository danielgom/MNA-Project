// Package security handles all security concerns such as password encryption and JWT generation and validation.
package security

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const passwordHashDifficulty = 10

// Hash produces an encrypted string.
func Hash(str string) (string, error) {
	hStr, err := bcrypt.GenerateFromPassword([]byte(str), passwordHashDifficulty)
	if err != nil {
		return "", fmt.Errorf("could encrypt password %w", err)
	}
	return string(hStr), nil
}

// CheckHash checks whether the password provided is the correct one.
func CheckHash(str, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(str))
	return err == nil
}
