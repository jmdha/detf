package main

import (
	"log"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	pb "detf/api"
)

type server struct {
	pb.UnimplementedDETFServer
}

func (s *server) Stream(
	stream pb.DETF_StreamServer,
) error {
	for i := 0; i < 32; i++ {
		stream.Send(&pb.Match {
			ID: uuid.New().String(),
		})
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("%v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDETFServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr());
	if err := s.Serve(lis); err != nil {
		log.Fatalf("%v", err)
	}
}
