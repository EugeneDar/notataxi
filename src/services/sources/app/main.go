package main

import (
	"app/src/services/sources/grpc"

	"log"
	"net"

	"google.golang.org/grpc"
)

type Service struct {
	GRPCServer *grpc.Server
}

func NewService() (*Service, error) {
	gRPCServer := grpc.NewServer()

	if err := grpcsources.Register(gRPCServer); err != nil {
		return nil, err
	}

	return &Service{
		GRPCServer: gRPCServer,
	}, nil
}

func main() {
	service, err := NewService()
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Source service started")

	log.Fatal(service.GRPCServer.Serve(l))
}
