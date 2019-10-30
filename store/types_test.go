package main

import (
	"testing"

	mocks "acbot/store/mocks"

	"github.com/golang/mock/gomock"
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

func TestInsertActivation(t *testing.T) {
	mr := getRepository(t)
	insertId, err := mr.InsertActivation(&testAct)
	assert.NoError(t, err, "Can't insertActivation!")
	getResult, err := mr.GetActivations(bson.D{primitive.E{Key: "_id", Value: insertId}})
	assert.NoError(t, err, "Can't get activations from Mongo!")
	assert.Equal(t, 1, len(getResult), "Insert activation error!")
	err = mr.Mongo.Database(mr.DbName).Collection(mr.ActivationCollection).Drop()
	assert.NoError(t, err, "Can't drop test collection!")
}

func TestMockInsertActivation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mongoMock := mocks.NewMockMongoInterface(mockCtrl)
	mongoMock.EXPECT().InsertActivation(&testAct).Return(testAct.ID.String(), nil)
	_, err := mongoMock.InsertActivation(&testAct)
	assert.NoError(t, err, "Can't insert activation!")
}
