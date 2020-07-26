package persistence

import (
	"trackpump/domain/model"
	"trackpump/domain/repository"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	companiesCollection = "users"
)

type mongoRepository struct {
	database *mgo.Database
}

// NewMongoRepository returns a new User Repository
func NewMongoRepository(database *mgo.Database) repository.UserRepository {
	ur := &mongoRepository{
		database: database,
	}
	ur.database.C(companiesCollection).EnsureIndex(mgo.Index{Key: []string{"email"}, Unique: true})
	return ur
}

func (m *mongoRepository) Save(company *model.User) (*model.User, error) {
	if err := m.database.C(companiesCollection).UpdateId(company.ID, company); err != nil {
		return company, m.database.C(companiesCollection).Insert(company)
	}
	return company, nil
}

func (m *mongoRepository) FindByEmail(email string) (*model.User, error) {
	var company model.User
	err := m.database.C(companiesCollection).Find(bson.M{"email": email}).One(&company)
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (m *mongoRepository) FindByPasswordResetToken(email string) (*model.User, error) {
	var company model.User
	err := m.database.C(companiesCollection).Find(bson.M{"passwordResetToken": email}).One(&company)
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (m *mongoRepository) FindByID(id string) (*model.User, error) {
	var company model.User
	err := m.database.C(companiesCollection).Find(bson.M{"_id": id}).One(&company)
	if err != nil {
		return nil, err
	}
	return &company, nil
}
