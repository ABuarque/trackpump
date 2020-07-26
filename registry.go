package main

import (
	"trackpump/adapter/controller"
	"trackpump/adapter/filestorage"
	"trackpump/adapter/id"
	"trackpump/adapter/password"
	"trackpump/adapter/persistence"
	"trackpump/auth"
	"trackpump/domain/repository"
	"trackpump/storage"
	"trackpump/usecase"
	"trackpump/usecase/service"

	"gopkg.in/mgo.v2"
)

type registry struct {
	client        *mgo.Database
	storageClient *storage.PCloudClient
	repository    repository.UserRepository
}

// Registry is an interface
type Registry interface {
	NewAppController() controller.AppController

	getRepository() repository.UserRepository
}

// NewRegistry returns a new registry
func NewRegistry(client *mgo.Database, storageClient *storage.PCloudClient) Registry {
	var repository repository.UserRepository
	if client == nil {
		repository = persistence.NewInMemoryRepository()
	} else {
		repository = persistence.NewMongoRepository(client)
	}
	return &registry{
		client:        client,
		storageClient: storageClient,
		repository:    repository,
	}
}

// injecting company repository
func (r *registry) getRepository() repository.UserRepository {
	return r.repository
}

// injecting id service
func (r *registry) getIDService() service.IDService {
	return id.New()
}

// injecting password service
func (r *registry) getPasswordService() service.PasswordService {
	return password.New()
}

// injecting auth service
func (r *registry) getAuthService() *auth.Auth {
	return auth.New()
}

// injecting storage service
func (r *registry) getStorageService() service.Storage {
	return filestorage.NewPcloudStorage(r.storageClient)
}

// injecting company use cases
func (r *registry) newCompanyUseCases() usecase.UseCases {
	return usecase.New(r.getRepository(), r.getPasswordService(), r.getIDService(), r.getStorageService())
}

// injecting customer controller
func (r *registry) NewAppController() controller.AppController {
	return controller.NewUsersController(r.newCompanyUseCases(), r.getAuthService())
}
