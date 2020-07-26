package persistence

import (
	"trackpump/domain/model"
	"trackpump/domain/repository"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	usersCollection = "users"
)

type mongoRepository struct {
	database *mgo.Database
}

// NewMongoRepository returns a new User Repository
func NewMongoRepository(database *mgo.Database) repository.UserRepository {
	ur := &mongoRepository{
		database: database,
	}
	ur.database.C(usersCollection).EnsureIndex(mgo.Index{Key: []string{"email"}, Unique: true})
	return ur
}

func (m *mongoRepository) Save(user *model.User) (*model.User, error) {
	if err := m.database.C(usersCollection).UpdateId(user.ID, user); err != nil {
		return user, m.database.C(usersCollection).Insert(user)
	}
	return user, nil
}

func (m *mongoRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := m.database.C(usersCollection).Find(bson.M{"email": email}).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *mongoRepository) FindByPasswordResetToken(email string) (*model.User, error) {
	var user model.User
	err := m.database.C(usersCollection).Find(bson.M{"passwordResetToken": email}).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *mongoRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	err := m.database.C(usersCollection).Find(bson.M{"_id": id}).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *mongoRepository) FindAll() ([]*model.User, error) {
	var users []*model.User
	err := m.database.C(usersCollection).Find(nil).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
