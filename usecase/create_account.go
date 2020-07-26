package usecase

import (
	"fmt"
	"strings"
	"time"
	"trackpump/domain/model"
	"trackpump/domain/repository"
	"trackpump/usecase/exception"
	"trackpump/usecase/service"

	"gopkg.in/validator.v2"
)

// CreateAccountInput is the use case input
type CreateAccountInput struct {
	Name     string `json:"name" validate:"max=255,min=2"`
	Email    string `json:"email" validate:"regexp=[A-Za-z0-9\\._-]+@[A-Za-z0-9]+\\..(\\.[A-Za-z]+)*"`
	Password string `json:"password" validate:"min=6"`
	Gender   int    `json:"gender" validate:"min=0,max=1"`
	Birth    string `json:"birth"`
	Height   int    `json:"height"`
}

// CreateAccountOutput is the use case output
type CreateAccountOutput struct {
	ID    string `json:"id"`
	Name  string `json:"name" validate:"max=255,min=2"`
	Email string `json:"email" validate:"regexp=[A-Za-z0-9\\._-]+@[A-Za-z0-9]+\\..(\\.[A-Za-z]+)*"`
}

type createAccount struct {
	repository      repository.UserRepository
	passwordService service.PasswordService
	idService       service.IDService
}

type createAccountUseCase interface {
	create(input *CreateAccountInput) (*CreateAccountOutput, error)
}

func newCreateAccountUseCase(repository repository.UserRepository, passwordService service.PasswordService, idService service.IDService) createAccountUseCase {
	return &createAccount{
		repository:      repository,
		passwordService: passwordService,
		idService:       idService,
	}
}

func (ca *createAccount) create(input *CreateAccountInput) (*CreateAccountOutput, error) {
	err := validator.Validate(input)
	if err != nil {
		return nil, exception.New(exception.InvalidParameters, err.Error(), err)
	}
	if _, err = ca.repository.FindByEmail(input.Email); err == nil {
		return nil, exception.New(exception.Conflict, fmt.Sprintf("email %s already in use", input.Email), nil)
	}
	password, err := ca.passwordService.Encrypt(input.Password)
	if err != nil {
		return nil, exception.New(exception.ProcessmentError, "failure to create user password", err)
	}
	id, err := ca.idService.Get()
	if err != nil {
		return nil, exception.New(exception.ProcessmentError, "failed to generate user id", err)
	}
	birth, err := time.Parse("2006-01-02", input.Birth)
	if err != nil {
		return nil, exception.New(exception.ProcessmentError, "failed to parse date", err)
	}
	user := model.User{
		ID:        id,
		Email:     strings.ToLower(input.Email),
		Password:  password,
		Name:      input.Name,
		Gender:    input.Gender,
		Height:    input.Height,
		Birth:     birth,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	d, err := ca.repository.Save(&user)
	if err != nil {
		return nil, exception.New(exception.Conflict, err.Error(), err)
	}
	return &CreateAccountOutput{
		ID:    d.ID,
		Email: d.Email,
		Name:  d.Name,
	}, nil
}
