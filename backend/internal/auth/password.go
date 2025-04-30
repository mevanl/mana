package auth

import "golang.org/x/crypto/bcrypt"

// Hashes password and returns string of hashed password
func HashPassword(password string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPasswordBytes), err
}

// Checks if passwords are same, returns bool
func CheckPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
