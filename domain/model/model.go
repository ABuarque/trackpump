package model

import "time"

// User is the user struct
type User struct {
	ID                 string             `bson:"_id"`
	Email              string             `bson:"email"`
	Name               string             `bson:"name"`
	Password           string             `bson:"password"`
	PasswordResetToken string             `bson:"passwordResetToken,omitempty"`
	Gender             int                `bson:"gender"`
	Birth              time.Time          `bson:"birth"`
	CreatedAt          time.Time          `bson:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt"`
	Height             int                `bson:"height"`
	Measurements       []*BodyMeasurement `bson:"measurements"`
}

// BodyMeasurement is data collected on a measurement
type BodyMeasurement struct {
	ID                     string    `bson:"id"`
	IssuedAt               time.Time `bson:"issuedAt"`
	Weight                 int       `bson:"weight"`
	AbdominalCircunference int       `bson:"abdominalCircunference"`
	Arm                    int       `bson:"arm"`
	Forearm                int       `bson:"forearm"`
	Calf                   int       `bson:"calf"`
	Neck                   int       `bson:"neck"`
	Hip                    int       `bson:"hip"`
	Thigh                  int       `bson:"thigh"`
	FrontalPicture         string    `bson:"frontalPicture"`
	SidePicture            string    `bson:"sidePicture"`
}
