package id

import (
	"trackpump/usecase/service"

	"github.com/google/uuid"
)

type id struct {
}

// New returns a new IDService implementation
func New() service.IDService {
	return &id{}
}

func (id *id) Get() (string, error) {
	return uuid.New().String(), nil
}
