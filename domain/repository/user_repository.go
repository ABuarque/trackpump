package repository

import "trackpump/domain/model"

// UserRepository defines how application will interact with db
type UserRepository interface {
	FindByID(id string) (*model.User, error)

	FindByEmail(email string) (*model.User, error)

	FindByPasswordResetToken(token string) (*model.User, error)

	Save(u *model.User) (*model.User, error)

	FindAll() ([]*model.User, error)

	SaveMeasurement(measurement *model.BodyMeasurement) (*model.BodyMeasurement, error)

	// The returned list containes only two elements whose are the last and
	// last but one measurements (ON THAT ORDER!)
	FindLastTwoMeasurements(userID string) ([]*model.BodyMeasurement, error)

	// It returns a list sorted by -issuedAt
	FindMeasurementsForProfile(userID string) ([]*model.BodyMeasurement, error)
}
