package main

import (
	"context"
	"grpc-lesson/gen/pb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func runCall(c pb.CallClient, in *pb.CallRequest) {
	log.Println("--- Unary ---")
	log.Printf("request: %#v\n", in)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := c.Call(ctx, in)
	if err != nil {
		log.Fatalf("could not call: %v", err)
	}
	log.Printf("response:%#v\n", resp)
}

func runBulkCall(c pb.CallClient, names []string) error {
	log.Println("--- ClientStreaming ---")
	stream, err := c.BulkCall(context.Background())
	if err != nil {
		return err
	}
	for _, name := range names {
		in := &pb.CallRequest{Name: name}
		log.Printf("request: %#v\n", in)
		if err := stream.Send(in); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		time.Sleep(time.Second * 1)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("response: %#v", res)
	return nil
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}
	defer conn.Close()
	c := pb.NewCallClient(conn)
	runCall(c, &pb.CallRequest{Name: "John"})
	runBulkCall(c, []string{"John", "Paul", "George", "Ringo"})
}
