package main

// go test -coverprofile=cover.out ./store && go tool cover -html=cover.out -o cover.html
import (
	"acbot/types"
	"context"
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	Client *mongo.Client
}

type MongoDatabase struct {
	Database *mongo.Database
}

type MongoCollection struct {
	Collection *mongo.Collection
}

// Validate envs for empties
func checkNoEmptyEnvs(envs []string) bool {
	for _, value := range envs {
		if value == "" {
			return false
		}
	}
	return true
}

// Get redis options from .env file
func GetDbUri(envFile string) (string, error) {
	var err error
	if envFile == "" {
		err = godotenv.Load()
	} else {
		err = godotenv.Load(envFile)
	}
	if err != nil {
		return "", err
	}
	addr := os.Getenv("MONGO_ADDR")
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASS")

	if !checkNoEmptyEnvs([]string{
		addr,
		user,
		password,
	}) {
		return "", errors.New("Your .env file have empty connect variables!")
	}
	return strings.Join([]string{
		"mongodb://",
		user,
		":",
		password,
		"@",
		addr,
	}, ""), nil
}

func (mc *MongoClient) Database(dbName string) *MongoDatabase {
	db := mc.Client.Database(dbName)
	return &MongoDatabase{Database: db}
}

func (md *MongoDatabase) Collection(colName string) *MongoCollection {
	collection := md.Database.Collection(colName)
	return &MongoCollection{Collection: collection}
}

func (mc *MongoClient) Connect(connectionUri string) (err error) {
	clientOptions := options.Client().ApplyURI(connectionUri)
	mc.Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	err = mc.Client.Ping(context.Background(), nil)
	return err
}

func (mc *MongoCollection) Insert(document interface{}) (*primitive.ObjectID, error) {
	insertResult, err := mc.Collection.InsertOne(context.Background(), document)
	if err != nil {
		return nil, err
	}
	objectInsertId, ok := insertResult.InsertedID.(primitive.ObjectID)
	if false == ok {
		return nil, errors.New("Can't parse InsertId to ObjectId!")
	}
	return &objectInsertId, err
}

func (mc *MongoCollection) GetByFilter(filter interface{}) (*mongo.Cursor, error) {
	options := options.Find()
	data, err := mc.Collection.Find(context.Background(), filter, options)
	return data, err
}

func (mc *MongoCollection) GetActivationsByFilter(filter interface{}) ([]*types.Activation, error) {
	data, err := mc.GetByFilter(filter)
	if err != nil {
		return nil, err
	}
	var results []*types.Activation
	for data.Next(context.Background()) {
		var elem *types.Activation
		err := data.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, elem)
	}
	return results, err
}

func (mc *MongoCollection) GetUsersByFilter(filter interface{}) ([]*types.User, error) {
	data, err := mc.GetByFilter(filter)
	if err != nil {
		return nil, err
	}
	var results []*types.User
	for data.Next(context.Background()) {
		var elem *types.User
		err := data.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, elem)
	}
	return results, err
}

func (mc *MongoCollection) Drop() error {
	return mc.Collection.Drop(context.Background())
}
