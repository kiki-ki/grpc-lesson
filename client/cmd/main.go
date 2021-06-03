package main

import (
	"context"
	"grpc-lesson/gen/pb"
	"log"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}
	defer conn.Close()
	c := pb.NewCallClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := c.Call(ctx, &pb.CallRequest{Name: "John"})
	if err != nil {
		log.Fatalf("could not call: %v", err)
	}

	log.Println(resp.GetMessage())
}
