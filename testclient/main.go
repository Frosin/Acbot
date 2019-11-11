package main

import (
	pb "acbot/proto/mongo"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

// func activationToPb(act *types.Activation) *pb.Activation {
// 	return &pb.Activation{
// 		User:      act.User,
// 		Activator: act.Activator,
// 		Complete:  act.Complete,
// 		Retry:     act.Retry,
// 	}
// }

func main() {
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())

	var testAct = &pb.Activation{
		ID:        "",
		Timestamp: time.Now().Format(time.RFC3339),
		User:      777777,
		Activator: 888888,
		Complete:  false,
		Retry:     false,
	}

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewMongoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.InsertActivation(ctx, testAct)
	if err != nil {
		log.Fatalf("Error by insert!", err)
	}
	fmt.Println("Inserted=", r.InsertId)

	filter := primitive.M{
		"complete": false,
	}
	byteFilter, err := json.Marshal(filter)
	gr, err := c.GetActivations(context.Background(), &pb.Filter{
		Value: string(byteFilter),
	})
	fmt.Printf("Get result=%v type=%T\n", len(gr.GetActivations()), gr)
	for i, v := range gr.GetActivations() {
		fmt.Printf("%v. Type=%T, value=%v\n", i, v, v)
	}
}
