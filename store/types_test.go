package main

// integration tests
import (
	pb "acbot/proto/mongo"
	"acbot/types"
	"context"
	"encoding/json"
	"log"
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
	//
	testAct.User = 777
	//
	log.Println("testAct=", testAct)
	insertId, err := mr.InsertActivation(context.Background(), types.GetActivationMy2Proto(testAct))
	assert.NoError(t, err, "Can't insertActivation!")
	//filter := strings.Join([]string{`{"_id": "`, insertId.GetInsertId(), `"}`}, "")

	filter := primitive.M{"_id": insertId.GetInsertId()}
	//filter := primitive.M{"user": 777}
	byteFilter, err := json.Marshal(filter)
	//
	//log.Println("filter=", string(byteFilter), "insertedId=", insertId.GetInsertId(), "len1=", len(string(byteFilter)), "len2=", len(strings.ReplaceAll(string(byteFilter), "\\", "")))
	//log.Println("string=", strings.ReplaceAll(string(byteFilter), "\\", ""))
	//
	assert.NoError(t, err, "Cant't marshal filter! Error=", err)
	getResult, err := mr.GetActivations(context.Background(), &pb.Filter{
		Value: string(byteFilter), //strings.ReplaceAll(string(byteFilter), "\\", ""),
	})
	//
	log.Println("string filter=", string(byteFilter), len(string(byteFilter)))
	log.Println("result ins=", insertId.GetInsertId())
	if len(getResult.GetActivations()) > 0 {
		log.Println("result mon=", getResult.GetActivations()[0].GetID())
	}
	log.Println("filter=", filter)

	//
	assert.NoError(t, err, "Can't get activations from Mongo!")
	assert.Equal(t, 1, len(getResult.GetActivations()), "Insert activation error!")
	err = mr.Mongo.Database(mr.DbName).Collection(mr.ActivationCollection).Drop()
	assert.NoError(t, err, "Can't drop test collection!")
}

func TestInsertGetUser(t *testing.T) {
	mr := getRepository(t)
	insertId, err := mr.InsertUser(context.Background(), types.GetUserMy2Proto(testUser))
	assert.NoError(t, err, "Can't insertUser!")
	byteFilter, err := json.Marshal(bson.M{"_id": insertId.GetInsertId()})
	assert.Error(t, err, "Cant't marshal filter!")
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
