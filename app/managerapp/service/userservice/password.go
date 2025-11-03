package userservice

import "golang.org/x/crypto/bcrypt"

const (
	bcryptCostAdjustment = 4
)

func VerifyPassword(hashedPassword, password string) bool {
	bErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return bErr == nil
}

func HashPassword(password string) (string, error) {
	bytes, hErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost+bcryptCostAdjustment)

	return string(bytes), hErr
}
