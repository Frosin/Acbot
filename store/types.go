package main

import (
	pb "acbot/proto/mongo"
	"acbot/types"
	"context"
	"encoding/json"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source=../store/types.go -package=mocks -destination=../store/mocks/mock_mongo.go MongoInterface

// type MongoInterface interface {
// 	Connect(uri string) (err error)
// 	InsertActivation(activation *types.Activation) (string, error)
// 	GetActivations(filter interface{}) ([]*types.Activation, error)
// 	InsertUser(user *types.User) (string, error)
// 	GetUsers(filter interface{}) ([]*types.User, error)
// }

type mongoRepository struct {
	pb.UnimplementedMongoServer
	Mongo                MongoClient
	Connected            bool
	UserCollection       string
	ActivationCollection string
	DbName               string
}

func (a *mongoRepository) Connect(uri string) (err error) {
	if uri == "" {
		uri, err = GetDbUri("")
		a.UserCollection = os.Getenv("MONGO_USER_COLLECTION")
		a.ActivationCollection = os.Getenv("MONGO_ACTIVATION_COLLECTION")
		a.DbName = os.Getenv("MONGO_DB")
	}
	err = a.Mongo.Connect(uri)
	a.Connected = true
	return
}

func (a *mongoRepository) InsertActivation(ctx context.Context, req *pb.Activation) (*pb.InsertResult, error) {
	activation, err := types.GetActivationProto2My(req)
	if err != nil {
		return &pb.InsertResult{}, err
	}
	activation.ID = primitive.NewObjectID().Hex()
	mongoResult, err := a.Mongo.Database(a.DbName).Collection(a.ActivationCollection).Insert(activation)
	return &pb.InsertResult{
		InsertId: mongoResult,
	}, err
}

func (a *mongoRepository) InsertUser(ctx context.Context, req *pb.User) (*pb.InsertResult, error) {
	user := types.GetUserProto2My(req)
	user.ID = primitive.NewObjectID().Hex()
	mongoResult, err := a.Mongo.Database(a.DbName).Collection(a.UserCollection).Insert(user)
	return &pb.InsertResult{
		InsertId: mongoResult,
	}, err
}

func (a *mongoRepository) GetActivations(ctx context.Context, req *pb.Filter) (*pb.GetActivationsResult, error) {
	var filter primitive.M
	err := json.Unmarshal([]byte(req.GetValue()), &filter)
	if err != nil {
		return &pb.GetActivationsResult{}, err
	}
	getResult, err := a.Mongo.Database(a.DbName).Collection(a.ActivationCollection).GetActivationsByFilter(filter)
	if err != nil {
		return &pb.GetActivationsResult{}, err
	}
	var protoResult []*pb.Activation
	for _, activation := range getResult {
		protoResult = append(protoResult, types.GetActivationMy2Proto(activation))
	}
	return &pb.GetActivationsResult{
		Activations: protoResult,
	}, nil
}

func (a *mongoRepository) GetUsers(ctx context.Context, req *pb.Filter) (*pb.GetUsersResult, error) {
	var filter primitive.M
	err := json.Unmarshal([]byte(req.GetValue()), &filter)
	if err != nil {
		return &pb.GetUsersResult{}, err
	}
	getResult, err := a.Mongo.Database(a.DbName).Collection(a.UserCollection).GetUsersByFilter(filter)
	if err != nil {
		return &pb.GetUsersResult{}, err
	}
	var protoResult []*pb.User
	for _, user := range getResult {
		protoResult = append(protoResult, types.GetUserMy2Proto(user))
	}
	return &pb.GetUsersResult{
		Users: protoResult,
	}, nil
}

/*
func (a *mongoRepository) InsertActivation(activation *types.Activation) (*primitive.ObjectID, error) {
	a.checkConnect()
	activation.ID = primitive.NewObjectID()
	return a.Mongo.Database(a.DbName).Collection(a.ActivationCollection).Insert(activation)
}


func (a *mongoRepository) GetActivations(filter interface{}) ([]*types.Activation, error) {
	a.checkConnect()
	return a.Mongo.Database(a.DbName).Collection(a.ActivationCollection).GetActivationsByFilter(filter)
}

func (a *mongoRepository) InsertUser(user *types.User) (*primitive.ObjectID, error) {
	a.checkConnect()
	user.ID = primitive.NewObjectID()
	return a.Mongo.Database(a.DbName).Collection(a.UserCollection).Insert(user)
}

func (a *mongoRepository) GetUsers(filter interface{}) ([]*types.User, error) {
	a.checkConnect()
	return a.Mongo.Database(a.DbName).Collection(a.UserCollection).GetUsersByFilter(filter)
}
*/
