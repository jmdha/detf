package main

import (
	"log"
	"context"
	"time"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "detf/api"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer conn.Close()
	client := pb.NewDETFClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	stream, sErr := client.Stream(ctx)
	if sErr != nil {
		log.Fatalf("%v", sErr)
	}
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("%v", err)
		}
		stream.Send(&pb.Result {
			ID: in.ID,
		})
	}
}

