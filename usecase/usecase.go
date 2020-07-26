package usecase

import (
	"trackpump/domain/repository"
	"trackpump/usecase/service"
)

type useCases struct {
	repository                 repository.UserRepository
	createAccountUseCase       createAccountUseCase
	loginUseCase               loginUseCase
	registerMeasurementUseCase registerMeasurementUseCase
}

// UseCases defines the possible use cases
type UseCases interface {
	CreateAccount(input *CreateAccountInput) (*CreateAccountOutput, error)

	Login(input *LoginInput) (*LoginOutput, error)

	RegisterMeasurement(input *RegisterMeasurementInput) error
}

// New creates a new use case set
func New(repository repository.UserRepository, passwordService service.PasswordService, idService service.IDService, storageService service.Storage) UseCases {
	return &useCases{
		repository:                 repository,
		createAccountUseCase:       newCreateAccountUseCase(repository, passwordService, idService),
		loginUseCase:               newLoginUseCase(repository, passwordService),
		registerMeasurementUseCase: newRegisterMeasurementUseCase(repository, storageService, idService),
	}
}

func (u *useCases) CreateAccount(input *CreateAccountInput) (*CreateAccountOutput, error) {
	return u.createAccountUseCase.create(input)
}

func (u *useCases) Login(input *LoginInput) (*LoginOutput, error) {
	return u.loginUseCase.login(input)
}

func (u *useCases) RegisterMeasurement(input *RegisterMeasurementInput) error {
	return u.registerMeasurementUseCase.register(input)
}
