package main

import (
	mocks "acbot/store/mocks"
	"acbot/types"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMockInsertActivations(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mongoMock := mocks.NewMockMongoInterface(mockCtrl)
	mongoMock.EXPECT().InsertActivation(testAct).Return(testAct.ID.String(), nil)
	_, err := mongoMock.InsertActivation(testAct)
	assert.NoError(t, err, "Can't insert activation!")
}

func TestMockGetActivation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mongoMock := mocks.NewMockMongoInterface(mockCtrl)
	filter := bson.D{primitive.E{Key: "activator", Value: 9876543}}
	mongoMock.EXPECT().GetActivations(filter).Return([]*types.Activation{testAct}, nil)

	_, err := mongoMock.GetActivations(filter)
	assert.NoError(t, err, "Can't insert activation!")
}

func TestMockGetUsers(t *testing.T) {
	// TODO mock tests
}

func TestMockInsertUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mongoMock := mocks.NewMockMongoInterface(mockCtrl)
	mongoMock.EXPECT().InsertUser(testUser).Return(testUser.ID.String(), nil)
	_, err := mongoMock.InsertUser(testUser)
	assert.NoError(t, err, "Can't insert activation!")
}
