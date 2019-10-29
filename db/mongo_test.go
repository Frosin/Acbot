package main

import (
	"acbot/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var testAct = types.Activation{
	Timestamp: time.Now(),
	User:      123456,
	Activator: 9876543,
	Complete:  false,
	Retry:     false,
}

type TestConnectSettings struct {
	client         *MongoClient
	databaseName   string
	collectionName string
}

func getClientConnection(t *testing.T) (*MongoClient, *MongoCollection) {
	var mongoClient MongoClient
	err := mongoClient.Connect("mongodb://root:123456@localhost")
	if !assert.NoError(t, err, "Failed to connect to Mongo!") {
		assert.FailNow(t, "No connect, tests failed!")
	}
	return &mongoClient, mongoClient.Database("admin").Collection("TestCollection")
}

func dropCollection(mongoCollection *MongoCollection, t *testing.T) {
	err := mongoCollection.Drop()
	assert.NoError(t, err, "Failed to drop collection!")
}

func TestInsert(t *testing.T) {
	_, collection := getClientConnection(t)
	insertId, err := collection.Insert(testAct)
	assert.NoError(t, err, "Failed to insert data!")
	assert.NotEmpty(t, insertId, "Bad insertId!")
	dropCollection(collection, t)
}

func TestGetActivationsByFilter(t *testing.T) {
	var testAct1 = &types.Activation{
		Timestamp: time.Now(),
		User:      111111,
		Activator: 9876543,
		Complete:  false,
		Retry:     false,
	}
	var testAct2 = &types.Activation{
		Timestamp: time.Now(),
		User:      222222,
		Activator: 9876543,
		Complete:  false,
		Retry:     false,
	}
	var testAct3 = &types.Activation{
		Timestamp: time.Now(),
		User:      333333,
		Activator: 8245677,
		Complete:  false,
		Retry:     false,
	}
	_, collection := getClientConnection(t)
	_, err := collection.Insert(testAct1)
	assert.NoError(t, err, "Failed to insert data!")
	_, err = collection.Insert(testAct2)
	assert.NoError(t, err, "Failed to insert data!")
	_, err = collection.Insert(testAct3)

	assert.NoError(t, err, "Failed to insert data!")
	filter := bson.D{primitive.E{Key: "activator", Value: 9876543}}
	results, err := collection.GetActivationsByFilter(filter)
	assert.Equal(t, 2, len(results))
	dropCollection(collection, t)
}

func TestGetUsersByFilter(t *testing.T) {
	var testUser1 = &types.User{
		ChatId:       0,
		FirstName:    "Ivan",
		LastName:     "Klepikov",
		UserName:     "Klepik3",
		Role:         "user",
		Active:       true,
		DeactiveTime: 0,
	}
	var testUser2 = &types.User{
		ChatId:       0,
		FirstName:    "Ivan",
		LastName:     "Klepikov",
		UserName:     "Klepik3",
		Role:         "admin",
		Active:       true,
		DeactiveTime: 0,
	}
	var testUser3 = &types.User{
		ChatId:       0,
		FirstName:    "Ivan",
		LastName:     "Klepikov",
		UserName:     "Klepik3",
		Role:         "helper",
		Active:       true,
		DeactiveTime: 0,
	}
	_, collection := getClientConnection(t)
	_, err := collection.Insert(testUser1)
	assert.NoError(t, err, "Failed to insert data!")
	_, err = collection.Insert(testUser2)
	assert.NoError(t, err, "Failed to insert data!")
	_, err = collection.Insert(testUser3)
	assert.NoError(t, err, "Failed to insert data!")
	filter := bson.D{primitive.E{Key: "role", Value: "helper"}}
	results, err := collection.GetActivationsByFilter(filter)
	assert.NoError(t, err, "Failed to get activations from DB!")
	assert.Equal(t, 1, len(results))
	dropCollection(collection, t)
}

func TestEmptyGetResult(t *testing.T) {
	_, collection := getClientConnection(t)
	filter := bson.D{primitive.E{Key: "role", Value: "tsar"}}
	results, err := collection.GetActivationsByFilter(filter)
	assert.NoError(t, err, "Failed to get activations from DB!")
	assert.Equal(t, 0, len(results))
}
