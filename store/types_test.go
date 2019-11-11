package main

// integration tests
import (
	pb "acbot/proto/mongo"
	"acbot/types"
	"context"
	"encoding/json"
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
	insertId, err := mr.InsertActivation(context.Background(), types.GetActivationMy2Proto(testAct))
	assert.NoError(t, err, "Can't insertActivation!")
	filter := primitive.M{"_id": insertId.GetInsertId()}
	byteFilter, err := json.Marshal(filter)
	assert.NoError(t, err, "Cant't marshal filter! Error=", err)
	getResult, err := mr.GetActivations(context.Background(), &pb.Filter{
		Value: string(byteFilter),
	})
	assert.NoError(t, err, "Can't get activations from Mongo!")
	assert.Equal(t, 1, len(getResult.GetActivations()), "Insert activation error!")
	err = mr.Mongo.Database(mr.DbName).Collection(mr.ActivationCollection).Drop()
	assert.NoError(t, err, "Can't drop test collection!")
}

func TestInsertGetUser(t *testing.T) {
	mr := getRepository(t)
	insertId, err := mr.InsertUser(context.Background(), types.GetUserMy2Proto(testUser))
	assert.NoError(t, err, "Can't insertUser!")
	filter := primitive.M{"_id": insertId.GetInsertId()}
	byteFilter, err := json.Marshal(filter)
	assert.NoError(t, err, "Cant't marshal filter! Error=", err)
	getResult, err := mr.GetUsers(context.Background(), &pb.Filter{
		Value: string(byteFilter),
	})
	assert.NoError(t, err, "Can't get user from Mongo!")
	assert.Equal(t, 1, len(getResult.GetUsers()), "Insert user error!")
	err = mr.Mongo.Database(mr.DbName).Collection(mr.UserCollection).Drop()
	assert.NoError(t, err, "Can't drop test collection!")
}

func TestInsertWithoutConnect(t *testing.T) {
	var mRep = &mongoRepository{
		Connected: false,
	}
	assert.Panics(t, func() {
		mRep.InsertActivation(context.Background(), types.GetActivationMy2Proto(testAct))
	})
}

func TestConvertActivationProto2My(t *testing.T) {
	testActProto := &pb.Activation{
		ID:        "num _id",
		Timestamp: "bad timestamp!",
		User:      98765,
		Activator: 12345,
		Complete:  false,
		Retry:     false,
	}
	mr := getRepository(t)
	_, err := mr.InsertActivation(context.Background(), testActProto)
	assert.Error(t, err, "Failed to get timestamp parse error!")
}

func TestGetUnmarshalError(t *testing.T) {
	filter := "its bad filter value!"
	mr := getRepository(t)
	_, err := mr.GetActivations(context.Background(), &pb.Filter{
		Value: filter,
	})
	assert.Error(t, err, "Failed to get unmarshal filter error!")
	_, err = mr.GetUsers(context.Background(), &pb.Filter{
		Value: filter,
	})
	assert.Error(t, err, "Failed to get unmarshal filter error!")
}

func TestGetByFilterError(t *testing.T) {
	mr := getRepository(t)
	_, mCol := getClientConnection(t)
	var badAct = &BadTestStruct{
		ID:        "it's mongo _id",
		Activator: "it's bad activator's num!",
	}
	mCol.Insert(badAct)

	filter := bson.D{primitive.E{Key: "activator", Value: "it's bad activator's num!"}}
	byteFilter, err := json.Marshal(filter)
	assert.NoError(t, err, "Error in marshalling filter!")
	_, err = mr.GetActivations(context.Background(), &pb.Filter{
		Value: string(byteFilter),
	})
	assert.Error(t, err, "Failed to get GetByFilter error!")
}
