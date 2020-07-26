package usecase

import (
	"fmt"
	"trackpump/domain/repository"
	"trackpump/usecase/exception"
	"trackpump/usecase/service"

	"gopkg.in/validator.v2"
)

// LoginInput is the use case input
type LoginInput struct {
	Email    string `json:"email" validate:"regexp=[A-Za-z0-9\\._-]+@[A-Za-z0-9]+\\..(\\.[A-Za-z]+)*"`
	Password string `json:"password" validate:"min=6"`
}

// LoginOutput is the use case output
type LoginOutput struct {
	ID    string `json:"id"`
	Name  string `json:"name" validate:"max=255,min=2"`
	Email string `json:"email" validate:"regexp=[A-Za-z0-9\\._-]+@[A-Za-z0-9]+\\..(\\.[A-Za-z]+)*"`
}

type login struct {
	repository      repository.UserRepository
	passwordService service.PasswordService
}

type loginUseCase interface {
	login(input *LoginInput) (*LoginOutput, error)
}

func newLoginUseCase(repository repository.UserRepository, passwordService service.PasswordService) loginUseCase {
	return &login{
		repository:      repository,
		passwordService: passwordService,
	}
}

func (l *login) login(input *LoginInput) (*LoginOutput, error) {
	if err := validator.Validate(input); err != nil {
		return nil, exception.New(exception.InvalidParameters, err.Error(), err)
	}
	user, err := l.repository.FindByEmail(input.Email)
	if err != nil {
		return nil, exception.New(exception.NotFound, fmt.Sprintf("user not found with email %s ", input.Email), err)
	}
	ok, err := l.passwordService.IsValid(input.Password, user.Password)
	if !ok || err != nil {
		return nil, exception.New(exception.NotFound, "invalid login params", err)
	}
	return &LoginOutput{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
