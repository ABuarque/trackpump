package persistence

import (
	"context"
	"fmt"
	"trackpump/domain/model"
	"trackpump/domain/repository"

	"cloud.google.com/go/datastore"
)

const (
	usersCollection        = "users"
	measurementsCollection = "measuremnts"
)

type datastoreRepository struct {
	client *datastore.Client
}

// NewDatastoreRepository returns a repository for datastore
func NewDatastoreRepository(client *datastore.Client) repository.UserRepository {
	return &datastoreRepository{
		client: client,
	}
}

func (dr *datastoreRepository) FindByID(id string) (*model.User, error) {
	var entities []*model.User
	q := datastore.NewQuery(usersCollection).Filter("ID =", id).Limit(1)
	if _, err := dr.client.GetAll(context.Background(), q, &entities); err != nil {
		return nil, fmt.Errorf("failed to search for id %s on collection %s, error %q", id, usersCollection, err)
	}
	if len(entities) == 0 {
		return nil, fmt.Errorf("not found any user with id %s", id)
	}
	return entities[0], nil
}

func (dr *datastoreRepository) FindByPasswordResetToken(token string) (*model.User, error) {
	var entities []*model.User
	q := datastore.NewQuery(usersCollection).Filter("PasswordResetToken =", token).Limit(1)
	if _, err := dr.client.GetAll(context.Background(), q, &entities); err != nil {
		return nil, fmt.Errorf("failed to search for password reset token %s on collection %s, error %q", token, usersCollection, err)
	}
	if len(entities) == 0 {
		return nil, fmt.Errorf("não foi encontrato nenhum usuário com token %s", token)
	}
	return entities[0], nil
}

func (dr *datastoreRepository) Save(u *model.User) (*model.User, error) {
	userKey := datastore.NameKey(usersCollection, u.ID, nil)
	if _, err := dr.client.Put(context.Background(), userKey, u); err != nil {
		return nil, fmt.Errorf("failed to save user on db, error %q", err)
	}
	return u, nil
}

func (dr *datastoreRepository) FindByEmail(email string) (*model.User, error) {
	var entities []*model.User
	q := datastore.NewQuery(usersCollection).Filter("Email =", email).Limit(1)
	if _, err := dr.client.GetAll(context.Background(), q, &entities); err != nil {
		return nil, fmt.Errorf("failed to search for email %s on collection %s, error %q", email, usersCollection, err)
	}
	if len(entities) == 0 {
		return nil, fmt.Errorf("not found any user with email %s", email)
	}
	return entities[0], nil
}

func (dr *datastoreRepository) FindAll() ([]*model.User, error) {
	var entities []*model.User
	q := datastore.NewQuery(usersCollection)
	if _, err := dr.client.GetAll(context.Background(), q, &entities); err != nil {
		return nil, fmt.Errorf("failed to find all users from db on collection %s, error %q", usersCollection, err)
	}
	return entities, nil
}

func (dr *datastoreRepository) SaveMeasurement(measurement *model.BodyMeasurement) (*model.BodyMeasurement, error) {
	measurementKey := datastore.NameKey(measurementsCollection, measurement.ID, nil)
	if _, err := dr.client.Put(context.Background(), measurementKey, measurement); err != nil {
		return nil, fmt.Errorf("failed to save measurement on db, error %q", err)
	}
	return measurement, nil
}

func (dr *datastoreRepository) FindLastTwoMeasurements(userID string) ([]*model.BodyMeasurement, error) {
	var entities []*model.BodyMeasurement
	q := datastore.NewQuery(measurementsCollection).Filter("UserID=", userID).Order("-IssuedAt").Limit(2)
	if _, err := dr.client.GetAll(context.Background(), q, &entities); err != nil {
		return nil, fmt.Errorf("failedt to fetch last two measurements on collection %s, error %q", usersCollection, err)
	}
	return entities, nil
}
