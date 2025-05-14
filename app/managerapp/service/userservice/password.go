package userservice

import "golang.org/x/crypto/bcrypt"

func VerifyPassword(hashedPassword, password string) bool {
	bErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return bErr == nil
}

func HashPassword(password string) (string, error) {
	bytes, hErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost+4)

	return string(bytes), hErr
}
