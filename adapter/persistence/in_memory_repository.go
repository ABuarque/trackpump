package persistence

import (
	"fmt"
	"trackpump/domain/model"
	"trackpump/domain/repository"
)

type inMemoryRepository struct {
	db map[string]*model.User
}

// NewInMemoryRepository returns an in memory repository
func NewInMemoryRepository() repository.UserRepository {
	return &inMemoryRepository{
		db: make(map[string]*model.User),
	}
}

func (im *inMemoryRepository) FindByID(id string) (*model.User, error) {
	for _, d := range im.db {
		if d.ID == id {
			return d, nil
		}
	}
	return nil, fmt.Errorf("was not found any user with id %s", id)
}

func (im *inMemoryRepository) FindByPasswordResetToken(token string) (*model.User, error) {
	for _, d := range im.db {
		if d.PasswordResetToken == token {
			return d, nil
		}
	}
	return nil, fmt.Errorf("was not found any user with token %s", token)
}

func (im *inMemoryRepository) Save(d *model.User) (*model.User, error) {
	im.db[d.ID] = d
	return d, nil
}

func (im *inMemoryRepository) FindByEmail(email string) (*model.User, error) {
	for _, d := range im.db {
		if d.Email == email {
			return d, nil
		}
	}
	return nil, fmt.Errorf("was not found any user with email %s", email)
}

func (im *inMemoryRepository) FindAll() ([]*model.User, error) {
	var users []*model.User
	for _, user := range im.db {
		users = append(users, user)
	}
	return users, nil
}
