package repository

import "trackpump/domain/model"

// UserRepository defines how application will interact with db
type UserRepository interface {
	FindByID(id string) (*model.User, error)

	FindByEmail(email string) (*model.User, error)

	FindByPasswordResetToken(token string) (*model.User, error)

	Save(u *model.User) (*model.User, error)
}
