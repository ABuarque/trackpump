package service

// PasswordService defines how password services should work
type PasswordService interface {
	Encrypt(password string) (string, error)

	IsValid(givenPasswordToCheck, gottenFromDataSource string) (bool, error)
}
