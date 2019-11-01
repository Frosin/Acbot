package main

// integration tests
import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTestingConnect(t *testing.T) {
	var mr mongoRepository
	err := mr.Connect("mongodb://root:123456@localhost")
	assert.NoError(t, err, "Can't connect!")
}

func getRepository(t *testing.T) *mongoRepository {
	var mr mongoRepository
	err := mr.Connect("")
	if !assert.NoError(t, err, "Can't connect!") {
		assert.FailNow(t, "Can't connect to mongo!")
	}
	return &mr
}

func TestInsertGetActivation(t *testing.T) {
	mr := getRepository(t)
	insertId, err := mr.InsertActivation(testAct)
	assert.NoError(t, err, "Can't insertActivation!")
	getResult, err := mr.GetActivations(bson.D{primitive.E{Key: "_id", Value: insertId}})
	assert.NoError(t, err, "Can't get activations from Mongo!")
	assert.Equal(t, 1, len(getResult), "Insert activation error!")
	err = mr.Mongo.Database(mr.DbName).Collection(mr.ActivationCollection).Drop()
	assert.NoError(t, err, "Can't drop test collection!")
}

func TestInsertGetUser(t *testing.T) {
	mr := getRepository(t)
	insertId, err := mr.InsertUser(testUser)
	assert.NoError(t, err, "Can't insertUser!")
	getResult, err := mr.GetUsers(bson.D{primitive.E{Key: "_id", Value: insertId}})
	assert.NoError(t, err, "Can't get user from Mongo!")
	assert.Equal(t, 1, len(getResult), "Insert user error!")
	err = mr.Mongo.Database(mr.DbName).Collection(mr.UserCollection).Drop()
	assert.NoError(t, err, "Can't drop test collection!")
}

func TestInsertWithoutConnect(t *testing.T) {
	var mRep = &mongoRepository{
		Connected: false,
	}
	assert.Panics(t, func() {
		mRep.InsertActivation(testAct)
	})
}
