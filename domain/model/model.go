package model

import "time"

// User is the user struct
type User struct {
	ID                 string
	Email              string
	Name               string
	Password           string
	PasswordResetToken string
	Gender             int
	Birth              time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Height             int
}

// BodyMeasurement is data collected on a measurement
type BodyMeasurement struct {
	ID                     string
	UserID                 string
	IssuedAt               time.Time
	Weight                 float64 // in grams
	AbdominalCircunference float64 // in cm
	Arm                    float64 // in cm
	Forearm                float64 // in cm
	Calf                   float64 // in cm
	Neck                   float64 // in cm
	Hip                    float64 // in cm
	Thigh                  float64 // in cm
	FrontalPicture         string
	SidePicture            string
	BodyFatPercentage      float64 // in %
	BodyMassIndex          float64
}
