package usecase

import (
	"strings"
	"trackpump/domain/repository"
	"trackpump/usecase/exception"
)

// LoadProfileInput is the use case input
type LoadProfileInput struct {
	ID string
}

// LoadProfileOutput is the use case output
type LoadProfileOutput struct {
	Labels             []string
	BodyFatPercentages []float64
	BodyMassIndexes    []float64
}

type loadProfile struct {
	repository repository.UserRepository
}

type loadProfileUseCase interface {
	load(input *LoadProfileInput) (*LoadProfileOutput, error)
}

func newLoadProfileUseCase(repository repository.UserRepository) loadProfileUseCase {
	return &loadProfile{
		repository: repository,
	}
}

func (lp *loadProfile) load(input *LoadProfileInput) (*LoadProfileOutput, error) {
	measurements, err := lp.repository.FindMeasurementsForProfile(input.ID)
	if err != nil {
		return nil, exception.New(exception.NotFound, err.Error(), err)
	}
	var labels []string
	var bodyFatPercentages []float64
	var bodyMassIndexes []float64
	for _, m := range measurements {
		labels = append(labels, strings.Split(m.IssuedAt.String(), " ")[0])
		bodyFatPercentages = append(bodyFatPercentages, m.BodyFatPercentage)
		bodyMassIndexes = append(bodyMassIndexes, m.BodyMassIndex)
	}
	return &LoadProfileOutput{
		Labels:             labels,
		BodyFatPercentages: bodyFatPercentages,
		BodyMassIndexes:    bodyMassIndexes,
	}, nil
}
