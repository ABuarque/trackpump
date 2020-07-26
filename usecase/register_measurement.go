package usecase

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"
	"trackpump/domain/model"
	"trackpump/domain/repository"
	"trackpump/usecase/exception"
	"trackpump/usecase/service"
)

// RegisterMeasurementInput is the use case input
type RegisterMeasurementInput struct {
	ID                     string `json:"id"`
	Weight                 int    `json:"weight"`
	AbdominalCircunference int    `json:"abdominalCircunference"`
	Arm                    int    `json:"arm"`
	Forearm                int    `json:"forearm"`
	Calf                   int    `json:"calf"`
	Neck                   int    `json:"neck"`
	Hip                    int    `json:"hip"`
	Thigh                  int    `json:"thigh"`
	FrontalPicture         string `json:"frontalPicture"`
	SidePicture            string `json:"sidePicture"`
}

type registerMeasurement struct {
	repository repository.UserRepository
	storage    service.Storage
	idService  service.IDService
}

type registerMeasurementUseCase interface {
	register(input *RegisterMeasurementInput) error
}

func newRegisterMeasurementUseCase(repository repository.UserRepository, storage service.Storage, idService service.IDService) registerMeasurementUseCase {
	return &registerMeasurement{
		repository: repository,
		storage:    storage,
		idService:  idService,
	}
}

func (r *registerMeasurement) register(input *RegisterMeasurementInput) error {
	user, err := r.repository.FindByID(input.ID)
	if err != nil {
		return exception.New(exception.NotFound, fmt.Sprintf("user not found with id %s", input.ID), err)
	}
	now := time.Now()
	frontalPicturesBytes, err := base64.StdEncoding.DecodeString(input.FrontalPicture)
	if err != nil {
		return exception.New(exception.ProcessmentError, "failed to decode frontal image from base 64, err %s", err)
	}
	frontalPicturesBuffer := new(bytes.Buffer)
	if _, err := frontalPicturesBuffer.Write(frontalPicturesBytes); err != nil {
		return fmt.Errorf("failed to get frontal pictures bytes, err %q", err)
	}
	frontalPictureName := fmt.Sprintf("%s_frontal-picture_%s.png", user.ID, timeToString(now))
	frontalPictureURL, err := r.storage.Put(frontalPictureName, frontalPicturesBuffer)
	if err != nil {
		return fmt.Errorf("failed to send frontal picture to storage, erro %q", err)
	}
	sidePicturesBytes, err := base64.StdEncoding.DecodeString(input.SidePicture)
	if err != nil {
		return exception.New(exception.ProcessmentError, "failed to decode side image from base 64, err %s", err)
	}
	sidePicturesBuffer := new(bytes.Buffer)
	if _, err := sidePicturesBuffer.Write(sidePicturesBytes); err != nil {
		return fmt.Errorf("failed to get frontal pictures bytes, err %q", err)
	}
	sidePictureName := fmt.Sprintf("%s_side-picture_%s.png", user.ID, timeToString(now))
	sidePictureURL, err := r.storage.Put(sidePictureName, sidePicturesBuffer)
	if err != nil {
		return fmt.Errorf("failed to send frontal picture to storage, erro %q", err)
	}
	id, err := r.idService.Get()
	if err != nil {
		return exception.New(exception.ProcessmentError, "failed to generate user id", err)
	}
	bodyMeasurement := model.BodyMeasurement{
		ID:                     id,
		IssuedAt:               time.Now(),
		Weight:                 input.Weight,
		AbdominalCircunference: input.AbdominalCircunference,
		Arm:                    input.Arm,
		Forearm:                input.Forearm,
		Calf:                   input.Calf,
		Neck:                   input.Neck,
		Hip:                    input.Hip,
		Thigh:                  input.Thigh,
		FrontalPicture:         frontalPictureURL,
		SidePicture:            sidePictureURL,
	}
	user.Measurements = append(user.Measurements, &bodyMeasurement)
	if _, err := r.repository.Save(user); err != nil {
		return exception.New(exception.ProcessmentError, "failed to update user metrics", err)
	}
	return nil
}

func timeToString(t time.Time) string {
	return t.Format("2006-01-02T15:04:05")
}
