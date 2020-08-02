package model

import "time"

// User is the user struct
type User struct {
	ID                 string    `bson:"_id"`
	Email              string    `bson:"email"`
	Name               string    `bson:"name"`
	Password           string    `bson:"password"`
	PasswordResetToken string    `bson:"passwordResetToken,omitempty"`
	Gender             int       `bson:"gender"`
	Birth              time.Time `bson:"birth"`
	CreatedAt          time.Time `bson:"createdAt"`
	UpdatedAt          time.Time `bson:"updatedAt"`
	Height             int       `bson:"height"`
}

// BodyMeasurement is data collected on a measurement
type BodyMeasurement struct {
	ID                     string    `bson:"_id"`
	UserID                 string    `bson:"userID"`
	IssuedAt               time.Time `bson:"issuedAt"`
	Weight                 float64   `bson:"weight"`                 // in grams
	AbdominalCircunference float64   `bson:"abdominalCircunference"` // in cm
	Arm                    float64   `bson:"arm"`                    // in cm
	Forearm                float64   `bson:"forearm"`                // in cm
	Calf                   float64   `bson:"calf"`                   // in cm
	Neck                   float64   `bson:"neck"`                   // in cm
	Hip                    float64   `bson:"hip"`                    // in cm
	Thigh                  float64   `bson:"thigh"`                  // in cm
	FrontalPicture         string    `bson:"frontalPicture"`
	SidePicture            string    `bson:"sidePicture"`
}
