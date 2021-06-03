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

func runUnaryCall(c pb.CallClient, in *pb.CallRequest) error {
	log.Println("--- Unary ---")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := c.UnaryCall(ctx, in)
	if err != nil {
		return err
	}
	log.Printf("response: %s\n", resp.GetMessage())
	return nil
}

func runClientStreamingCall(c pb.CallClient, names []string) error {
	log.Println("--- ClientStreaming ---")
	stream, err := c.ClientStreamingCall(context.Background())
	if err != nil {
		return err
	}
	for _, name := range names {
		in := &pb.CallRequest{Name: name}
		if err := stream.Send(in); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		time.Sleep(time.Second)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("response: %s", res.GetMessage())
	return nil
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}
	defer conn.Close()
	c := pb.NewCallClient(conn)

	err = runUnaryCall(c, &pb.CallRequest{Name: "John"})
	if err != nil {
		log.Fatalln(err)
	}
	err = runClientStreamingCall(c, []string{"John", "Paul", "George", "Ringo"})
	if err != nil {
		log.Fatalln(err)
	}
}
