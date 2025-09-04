package hash

import "golang.org/x/crypto/bcrypt"

func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func CheckPassword(hashStr, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashStr), []byte(plain)) == nil
}
