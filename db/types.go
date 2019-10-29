package main

import (
	"acbot/types"
	"os"
)

type MongoActivation interface {
}

type MongoUser interface {
}

type mongoRepository struct {
	Mongo                MongoClient
	Connected            bool
	UserCollection       string
	ActivationCollection string
	DbName               string
}

func (a *mongoRepository) checkConnect() {
	if a.Connected == false {
		panic("No connection to database! Connect() not called!")
	}
}

func (a *mongoRepository) Connect(uri string) (err error) {
	if uri == "" {
		uri, err = GetDbUri()
		a.UserCollection = os.Getenv("MONGO_USER_COLLECTION")
		a.ActivationCollection = os.Getenv("MONGO_ACTIVATION_COLLECTION")
		a.DbName = os.Getenv("MONGO_DB")
	}
	err = a.Mongo.Connect(uri)
	a.Connected = true
	return
}

func (a *mongoRepository) InsertActivation(activation *types.Activation) (string, error) {
	a.checkConnect()
	return a.Mongo.Database(a.DbName).Collection(a.ActivationCollection).Insert(activation)
}

func (a *mongoRepository) GetActivations(filter interface{}) ([]*types.Activation, error) {
	a.checkConnect()
	return a.Mongo.Database(a.DbName).Collection(a.ActivationCollection).GetActivationsByFilter(filter)
}

func (a *mongoRepository) InsertUser(user *types.User) (string, error) {
	a.checkConnect()
	return a.Mongo.Database(a.DbName).Collection(a.UserCollection).Insert(user)
}

func (a *mongoRepository) GetUsers(filter interface{}) ([]*types.User, error) {
	a.checkConnect()
	return a.Mongo.Database(a.DbName).Collection(a.UserCollection).GetUsersByFilter(filter)
}
