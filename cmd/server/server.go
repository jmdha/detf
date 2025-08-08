package main

import (
	"log"
	"net"
	"io"
	"fmt"

	"google.golang.org/grpc"
	pb "detf/api"
)

type Server struct {
	pb.UnimplementedDETFServer
}

func SendMatch(
	stream pb.DETF_StreamServer,
) error {
	initial_match, err := NextMatch()
	if err != nil {
		return err
	}
	stream.Send(&pb.Match {
		ID: initial_match.id,
	})
	return nil
}

func (s *Server) Stream(
	stream pb.DETF_StreamServer,
) error {
	{
		err := SendMatch(stream)
		if err != nil {
			return err
		}
	}
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		HandleResult(
			result {
				id: in.ID,
			},
		)
		{
			err := SendMatch(stream)
			if err != nil {
				return err
			}
		}
	}
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
