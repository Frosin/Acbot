package main

import (
	pb "acbot/proto/mongo"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedMongoServer
	MongoRepository mongoRepository
}

func (s *Server) InsertActivation(ctx context.Context, req *pb.Activation) (*pb.ActivationInsertResult, error) {
	fmt.Println("InsertActivation called!", req)

	//s.MongoRepository.InsertActivation()
	return &pb.ActivationInsertResult{
		InsertId: "success",
	}, nil
}

func panicIfError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func main() {
	var mongoServer Server
	//err := mongoServer.MongoRepository.Connect("")
	err := mongoServer.MongoRepository.Connect("mongodb://root:123456@localhost")

	panicIfError("Failed to connect: ", err)
	fmt.Println("Starting server...")

	l, err := net.Listen("tcp", ":8081")
	panicIfError("Failed to listen", err)

	s := grpc.NewServer()
	pb.RegisterMongoServer(s, &Server{})

	err = s.Serve(l)
	panicIfError("Failed to serve", err)

}
