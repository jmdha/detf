package main

import (
	"log"
	"net"
	"fmt"
	"context"

	"google.golang.org/grpc"
	pb "detf/api"
)

type Server struct {
	pb.UnimplementedDETFServer
}

func (s *Server) RequestMatch(_ context.Context, in *pb.Empty) (*pb.Match, error) {
	match, err := NextMatch()
	if err != nil {
		return &pb.Match {}, err
	}
	return &match, nil
}

func (s *Server) SendResult(_ context.Context, in *pb.Result) (*pb.Empty, error) {
	HandleResult(*in)
	return &pb.Empty {}, nil
}

func ServerStart(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("%v", err)
	}
	server := grpc.NewServer()
	pb.RegisterDETFServer(server, &Server{})
	log.Printf("server listening at %v", lis.Addr());
	if err := server.Serve(lis); err != nil {
		log.Fatalf("%v", err)
	}
}
