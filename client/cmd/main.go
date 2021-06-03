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

func runUnaryCall(c pb.CallClient, name string) error {
	log.Println("--- Unary ---")
	in := &pb.CallRequest{Name: name}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	res, err := c.UnaryCall(ctx, in)
	if err != nil {
		return err
	}
	log.Printf("response: %s\n", res.GetMessage())
	return nil
}

func runClientStreamingCall(c pb.CallClient, names []string) error {
	log.Println("--- ClientStreaming ---")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	stream, err := c.ClientStreamingCall(ctx)
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

func runServerStreamingCall(c pb.CallClient, name string, responseCnt int32) error {
	log.Println("--- ServerStreaming ---")
	in := &pb.ServerStreamingCallRequest{Name: name, ResponseCnt: responseCnt}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	stream, err := c.ServerStreamingCall(ctx, in)
	if err != nil {
		return err
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("response: %s", res.GetMessage())
	}
	return nil
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}
	defer conn.Close()
	c := pb.NewCallClient(conn)

	if err = runUnaryCall(c, "John"); err != nil {
		log.Fatalln(err)
	}
	if err = runClientStreamingCall(c, []string{"John", "Paul", "George", "Ringo"}); err != nil {
		log.Fatalln(err)
	}
	if err = runServerStreamingCall(c, "John", 10); err != nil {
		log.Fatalln(err)
	}
}
