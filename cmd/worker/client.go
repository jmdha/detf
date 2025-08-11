package main

import (
	"log"
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "detf/api"
)

var client pb.DETFClient

func RequestMatch() (pb.Match, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	match, err := client.RequestMatch(ctx, &pb.Empty {})
	if err != nil {
		return pb.Match {}, err
	} else {
		return *match, err
	}
}

func SendResult(res pb.Result) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	_, err := client.SendResult(ctx, &res)
	if err != nil {
		log.Printf("Failed to send result with error: %v", err)
	}
}

func InitClient(ip string) {
	conn, err := grpc.NewClient(
		ip,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to server with error: %v", err)
	}
	client = pb.NewDETFClient(conn)
}
