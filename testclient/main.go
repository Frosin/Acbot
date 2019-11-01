package main

import (
	pb "acbot/proto/mongo"
	"acbot/types"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
)

func activationToPb(act *types.Activation) *pb.Activation {
	return &pb.Activation{
		User:      act.User,
		Activator: act.Activator,
		Complete:  act.Complete,
		Retry:     act.Retry,
	}
}

func main() {
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())

	var testAct = &types.Activation{
		// ID:        primitive.NewObjectID(),
		// Timestamp: time.Now(),
		User:      123456,
		Activator: 9876543,
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

	r, err := c.InsertActivation(ctx, activationToPb(testAct))
	if err != nil {
		log.Fatalf("Error by insert!", err)
	}
	fmt.Println(r.InsertId)
}
