package db

import "gopkg.in/mgo.v2"

// New returns a new db instance
func New(url, database string) (*mgo.Database, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	db := session.DB(database)
	return db, nil
}
