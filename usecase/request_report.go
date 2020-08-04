package usecase

import (
	"fmt"
	"log"
	"strings"
	"time"
	"trackpump/domain/model"
	"trackpump/domain/repository"
	"trackpump/usecase/exception"
	"trackpump/usecase/service"
)

const (
	male   = 1
	female = 0
)

type requestReport struct {
	repository   repository.UserRepository
	notification service.Notification
}

type requestReportUseCase interface {
	process() error
}

func newWeeklyWorkoutReport(repository repository.UserRepository, notification service.Notification) requestReportUseCase {
	return &requestReport{
		repository:   repository,
		notification: notification,
	}
}

func (rr *requestReport) process() error {
	log.Println("starting weekly workout report")
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
			payload := service.WeeklyReportPayload{
				Email:  user.Email,
				Name:   user.Name,
				Report: report,
			}
			if err := rr.notification.SendWeeklyReport(&payload); err != nil {
				log.Printf("failed to send report to email %s\n", user.Email)
			}
		} else {
			log.Printf("user %s does not have 2 measurements yet", user.ID)
		}
	}
	return nil
}

func getWorkoutReport(lastMeasure, lastButOneMeasure *model.BodyMeasurement, height, gender int, birth time.Time) (string, error) {
	var report strings.Builder
	weightDiff := (lastMeasure.Weight - lastButOneMeasure.Weight) / float64(1000) // converting to Kg
	if _, err := report.WriteString(fmt.Sprintf("Last weight: %fkg (diff: %fkg)\n", float64(lastMeasure.Weight)/1000.0, weightDiff)); err != nil {
		return "", fmt.Errorf("failed to write weight, erro %q", err)
	}
	abdominalCircunferenceDiff := lastMeasure.AbdominalCircunference - lastButOneMeasure.AbdominalCircunference // in cm
	if _, err := report.WriteString(fmt.Sprintf("Last abdominal circunference: %fcm (diff: %fcm)\n", lastMeasure.AbdominalCircunference, abdominalCircunferenceDiff)); err != nil {
		return "", fmt.Errorf("failed to write abdominal circunference. erro %q", err)
	}
	armDiff := lastMeasure.Arm - lastButOneMeasure.Arm // in cm
	if _, err := report.WriteString(fmt.Sprintf("Last arm measure: %fcm (diff: %fcm)\n", lastMeasure.Arm, armDiff)); err != nil {
		return "", fmt.Errorf("failed to write arm. erro %q", err)
	}
	forearmDiff := lastMeasure.Forearm - lastButOneMeasure.Forearm // in cm
	if _, err := report.WriteString(fmt.Sprintf("Last forearm measure: %fcm (diff: %fcm)\n", lastMeasure.Forearm, forearmDiff)); err != nil {
		return "", fmt.Errorf("failed to write forearm. erro %q", err)
	}
	calfDiff := lastMeasure.Calf - lastButOneMeasure.Calf // in cm
	if _, err := report.WriteString(fmt.Sprintf("Last calf measure: %fcm (diff: %fcm)\n", lastMeasure.Calf, calfDiff)); err != nil {
		return "", fmt.Errorf("failed to write calf. erro %q", err)
	}
	neckDiff := lastMeasure.Neck - lastButOneMeasure.Neck // in cm
	if _, err := report.WriteString(fmt.Sprintf("Last neck measure: %fcm (diff: %fcm)\n", lastMeasure.Neck, neckDiff)); err != nil {
		return "", fmt.Errorf("failed to write neck. erro %q", err)
	}
	hipDiff := lastMeasure.Hip - lastButOneMeasure.Hip // in cm
	if _, err := report.WriteString(fmt.Sprintf("Last hip measure: %fcm (diff: %fcm)\n", lastMeasure.Hip, hipDiff)); err != nil {
		return "", fmt.Errorf("failed to write hip. erro %q", err)
	}
	thighDiff := lastMeasure.Thigh - lastButOneMeasure.Thigh // in cm
	if _, err := report.WriteString(fmt.Sprintf("Last thigh measure: %fcm (diff: %fcm)\n", lastMeasure.Thigh, thighDiff)); err != nil {
		return "", fmt.Errorf("failed to write thight. erro %q", err)
	}
	bodyMassIndexDiff := lastMeasure.BodyMassIndex - lastButOneMeasure.BodyMassIndex
	if _, err := report.WriteString(fmt.Sprintf("BMI: %f [%s] (%f)\n", lastMeasure.BodyMassIndex, getBodyMassIndexStatus(lastMeasure.BodyMassIndex), bodyMassIndexDiff)); err != nil {
		return "", fmt.Errorf("failed to write body mass index, erro %q", err)
	}
	bodyFatPercentageDiff := lastMeasure.BodyFatPercentage - lastButOneMeasure.BodyFatPercentage
	if _, err := report.WriteString(fmt.Sprintf("Body fat percentage %f%s (%f%s) \n", lastMeasure.BodyFatPercentage, "%", bodyFatPercentageDiff, "%")); err != nil {
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
