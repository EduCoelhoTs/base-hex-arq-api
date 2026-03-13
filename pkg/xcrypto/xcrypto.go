package xcrypto

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytePass := []byte(password)
	hashedPass, err := bcrypt.GenerateFromPassword(bytePass, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}

func ComparePassword(hashedPassword, password string) error {
	byteHashedPass := []byte(hashedPassword)
	bytePass := []byte(password)

	return bcrypt.CompareHashAndPassword(byteHashedPass, bytePass)
}
