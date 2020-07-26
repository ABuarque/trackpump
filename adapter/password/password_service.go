package password

import (
	"fmt"
	"trackpump/usecase/service"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
}

// New retuns a new password service
func New() service.PasswordService {
	return &password{}
}

func (p *password) Encrypt(password string) (string, error) {
	stringToEncryptAsBytes := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(stringToEncryptAsBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password %q", err)
	}
	return string(hash), nil
}

func (p *password) IsValid(stringToCheck, stringFromDataSource string) (bool, error) {
	stringToCheckAsBytes := []byte(stringToCheck)
	stringFromDataSourceAsBytes := []byte(stringFromDataSource)
	err := bcrypt.CompareHashAndPassword(stringFromDataSourceAsBytes, stringToCheckAsBytes)
	if err != nil {
		return false, fmt.Errorf("failed to verify password hash: %q", err)
	}
	return true, nil
}
