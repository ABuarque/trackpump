package usecase

import (
	"trackpump/domain/repository"
	"trackpump/usecase/exception"
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
		if len(user.Measurements) >= 2 {
			// TODO sort measurements
			// get last two
			// make diff
		}
	}
	return nil
}
