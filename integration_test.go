package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"trackpump/domain/model"
	"trackpump/usecase"
	"trackpump/usecase/exception"

	"github.com/labstack/echo"
)

func TestCreateAccount(t *testing.T) {
	registry := NewRegistry(nil, nil, "EMAIL", "PASSWORD")
	controller := registry.NewAppController()
	inputUseCase := `{"name":"Aurelio Buarque", "email":"abuarquemf@gmail.com", "password":"123455678", "gender":0, "birth":"1997-11-29", "height":172}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(inputUseCase))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	controller.Create(c)
	outputUseCase := usecase.CreateAccountOutput{}
	if err := json.Unmarshal(rec.Body.Bytes(), &outputUseCase); err != nil {
		t.Errorf("expected error nil when unmarshaling response, erro %q", err)
	}
	if outputUseCase.Email != "abuarquemf@gmail.com" {
		t.Errorf("expected email to be abuarquemf@gmail.com, got %s", outputUseCase.Email)
	}
	if outputUseCase.Name != "Aurelio Buarque" {
		t.Errorf("expected email to be Aurelio Buarque, got %s", outputUseCase.Name)
	}
	if outputUseCase.ID == "" {
		t.Errorf("expected to have id different from an empty string")
	}
}

func TestCreateAccountWithEmailAlreadyOnDB(t *testing.T) {
	registry := NewRegistry(nil, nil, "EMAIL", "PASSWORD")
	controller := registry.NewAppController()
	registry.getRepository().Save(&model.User{
		Email: "abuarquemf@gmail.com",
	})
	inputUseCase := `{"name":"Aurelio Buarque", "email":"abuarquemf@gmail.com", "password":"123455678", "gender":0, "birth":"1997-11-29", "height":172}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(inputUseCase))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	controller.Create(c)
	errResponse := exception.Error{}
	if err := json.Unmarshal(rec.Body.Bytes(), &errResponse); err != nil {
		t.Errorf("expected error nil when unmarshaling response, erro %q", err)
	}
	if errResponse.Message != "email abuarquemf@gmail.com already in use" {
		t.Errorf("expected to get \"email abuarquemf@gmail.com already in use\", got %s", errResponse.Message)
	}
	if errResponse.Code != 409 {
		t.Errorf("expected code to be 409, got %d", errResponse.Code)
	}
}

func TestLoginWithEmailNotPresentOnDB(t *testing.T) {
	registry := NewRegistry(nil, nil, "EMAIL", "PASSWORD")
	controller := registry.NewAppController()
	inputUseCase := `{"email":"abuarquemf@gmail.com", "password":"12345678"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(inputUseCase))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	controller.Login(c)
	errResponse := exception.Error{}
	if err := json.Unmarshal(rec.Body.Bytes(), &errResponse); err != nil {
		t.Errorf("expected error nil when unmarshaling response, erro %q", err)
	}
	if errResponse.Message != "user not found with email abuarquemf@gmail.com " {
		t.Errorf("expeted message \"user not found with email abuarquemf@gmail.com\", got %s", errResponse.Message)
	}
	if errResponse.Code != 404 {
		t.Errorf("expected code 404, got %d", errResponse.Code)
	}
}

func TestLoginWithEmailPresentOnDB(t *testing.T) {
	registry := NewRegistry(nil, nil, "EMAIL", "PASSWORD")
	controller := registry.NewAppController()
	registry.getRepository().Save(&model.User{
		Email:    "abuarquemf@gmail.com",
		Name:     "Aurelio Buarque",
		ID:       "505",
		Password: "$2a$10$X4m8N.KozNblKzHwm.2KpudOdq5k0TyNvFqBzo/G23eDgN4HKtyB6",
	})
	inputUseCase := `{"email":"abuarquemf@gmail.com", "password":"1234567"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(inputUseCase))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	controller.Login(c)
	outputUseCase := usecase.LoginOutput{}
	if err := json.Unmarshal(rec.Body.Bytes(), &outputUseCase); err != nil {
		t.Errorf("expected error nil when unmarshaling response, erro %q", err)
	}
	if outputUseCase.Email != "abuarquemf@gmail.com" {
		t.Errorf("expected email abuarquemf@gmail.com, got %s", outputUseCase.Email)
	}
	if outputUseCase.Name != "Aurelio Buarque" {
		t.Errorf("expected name Aurelio Buarque, got %s", outputUseCase.Name)
	}
	if outputUseCase.ID != "505" {
		t.Errorf("exepcted ID 505, got %s", outputUseCase.ID)
	}
}
