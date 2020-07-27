package usecase

import (
	"fmt"
	"log"
	"strings"
	"time"
	"trackpump/domain/model"
	"trackpump/domain/repository"
	"trackpump/usecase/exception"
)

const (
	male   = 1
	female = 0
)

type requestReport struct {
	repository repository.UserRepository
}

type requestReportUseCase interface {
	process() error
}

func newWeeklyWorkoutReport(repository repository.UserRepository) requestReportUseCase {
	return &requestReport{
		repository: repository,
	}
}

func (rr *requestReport) process() error {
	users, err := rr.repository.FindAll()
	if err != nil {
		return exception.New(exception.ProcessmentError, "failed to retrieve all users from db", err)
	}
	for _, user := range users {
		lastMeasurements, err := rr.repository.FindLastTwoMeasurements(user.ID)
		if err != nil {
			log.Printf("error on process for user %s, erro %q", user.ID, err)
			continue
		}
		if len(lastMeasurements) >= 2 {
			lastMeasure := lastMeasurements[0]
			lastButOneMeasure := lastMeasurements[1]
			report, err := getWorkoutReport(lastMeasure, lastButOneMeasure, user.Height, user.Gender, user.Birth)
			if err != nil {
				return exception.New(exception.ProcessmentError, fmt.Sprintf("failed to build report, erro %q", err), err)
			}
			fmt.Println(report)
		}
		// TODO send report via email
	}
	return nil
}

func getWorkoutReport(lastMeasure, lastButOneMeasure *model.BodyMeasurement, height, gender int, birth time.Time) (string, error) {
	var report strings.Builder
	weightDiff := (lastMeasure.Weight - lastButOneMeasure.Weight) / 1000 // converting to Kg
	if _, err := report.WriteString(fmt.Sprintf("Last weight: %fkg (diff: %dkg)\n", float64(lastMeasure.Weight)/1000.0, weightDiff)); err != nil {
		return "", fmt.Errorf("failed to write weight, erro %q", err)
	}
	abdominalCircunferenceDiff := lastMeasure.AbdominalCircunference - lastButOneMeasure.AbdominalCircunference // in cm
	if _, err := report.WriteString(fmt.Sprintf("Last abdominal circunference: %dcm (diff: %dcm)\n", lastMeasure.AbdominalCircunference, abdominalCircunferenceDiff)); err != nil {
		return "", fmt.Errorf("failed to write abdominal circunference. erro %q", err)
	}
	armDiff := lastMeasure.Arm - lastButOneMeasure.Arm // in cm
	if _, err := report.WriteString(fmt.Sprintf("Last arm measure: %dcm (diff: %dcm)\n", lastMeasure.Arm, armDiff)); err != nil {
		return "", fmt.Errorf("failed to write arm. erro %q", err)
	}
	forearmDiff := lastMeasure.Forearm - lastButOneMeasure.Forearm // in cm
	if _, err := report.WriteString(fmt.Sprintf("Last forearm measure: %dcm (diff: %dcm)\n", lastMeasure.Forearm, forearmDiff)); err != nil {
		return "", fmt.Errorf("failed to write forearm. erro %q", err)
	}
	calfDiff := lastMeasure.Calf - lastButOneMeasure.Calf // in cm
	if _, err := report.WriteString(fmt.Sprintf("Last calf measure: %dcm (diff: %dcm)\n", lastMeasure.Calf, calfDiff)); err != nil {
		return "", fmt.Errorf("failed to write calf. erro %q", err)
	}
	neckDiff := lastMeasure.Neck - lastButOneMeasure.Neck // in cm
	if _, err := report.WriteString(fmt.Sprintf("Last neck measure: %dcm (diff: %dcm)\n", lastMeasure.Neck, neckDiff)); err != nil {
		return "", fmt.Errorf("failed to write neck. erro %q", err)
	}
	hipDiff := lastMeasure.Hip - lastButOneMeasure.Hip // in cm
	if _, err := report.WriteString(fmt.Sprintf("Last hip measure: %dcm (diff: %dcm)\n", lastMeasure.Hip, hipDiff)); err != nil {
		return "", fmt.Errorf("failed to write hip. erro %q", err)
	}
	thighDiff := lastMeasure.Thigh - lastButOneMeasure.Thigh // in cm
	if _, err := report.WriteString(fmt.Sprintf("Last thigh measure: %dcm (diff: %dcm)\n", lastMeasure.Thigh, thighDiff)); err != nil {
		return "", fmt.Errorf("failed to write thight. erro %q", err)
	}
	bodyMassIndex := (float64(lastMeasure.Weight) / 1000.0) / (float64(height) * float64(height) / 10000.0) // kg/m^2
	if _, err := report.WriteString(fmt.Sprintf("BMI: %f (%s)\n", bodyMassIndex, getBodyMassIndexStatus(bodyMassIndex))); err != nil {
		return "", fmt.Errorf("failed to write body mass index, erro %q", err)
	}
	bodyFatPercentage := getBodyFatPercentage(bodyMassIndex, gender, birth)
	if _, err := report.WriteString(fmt.Sprintf("Body fat percentage %f%s \n", bodyFatPercentage, "%")); err != nil {
		return "", fmt.Errorf("failed to write body fat percentage, erro %q", err)
	}
	return report.String(), nil
}

func getBodyMassIndexStatus(bodyMassIndex float64) string {
	if bodyMassIndex < 18.5 {
		return "Thinness"
	} else if bodyMassIndex >= 18.5 && bodyMassIndex <= 24.9 {
		return "Regular"
	} else if bodyMassIndex >= 25 && bodyMassIndex <= 29.9 {
		return "Over weight"
	} else if bodyMassIndex >= 30 && bodyMassIndex <= 39.9 {
		return "Obesity"
	} else {
		return "Severe Obesity"
	}
}

func getBodyFatPercentage(bodyMassIndex float64, gender int, birth time.Time) float64 {
	yearOfBorn, _, _ := birth.Date()
	currentYear, _, _ := time.Now().Date()
	age := currentYear - yearOfBorn
	bodyFatPercentage := (1.2 * bodyMassIndex) + (0.23 * float64(age)) - (10.8 * float64(gender)) - 5.4 // male 1, female 0
	return bodyFatPercentage
}
