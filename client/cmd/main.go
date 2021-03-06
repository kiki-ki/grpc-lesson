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

func runServerStreamingCall(c pb.CallClient, name string, responseCnt uint32) error {
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

func runBidirectionalStreamingCall(c pb.CallClient, names []string) error {
	log.Println("--- BidirectionalStreaming ---")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	stream, err := c.BidirectionalStreamingCall(ctx)
	if err != nil {
		return err
	}
	done := make(chan struct{})
	go recv(done, stream)
	if err := send(names, stream); err != nil {
		return err
	}
	<- done
	return nil
}

func send(names []string, stream pb.Call_BidirectionalStreamingCallClient) error {
	for _, name := range names {
		in := &pb.CallRequest{Name:name}
		if err := stream.Send(in); err != nil {
			return err
		}
	}
	if err := stream.CloseSend(); err != nil {
		return err
	}
	return nil
}

func recv(done chan struct{}, stream pb.Call_BidirectionalStreamingCallClient) {
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			close(done)
			return
		}
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("response: %v\n", res.CallCounter)
	}
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
	names := []string{"John", "Paul", "George", "Ringo"}
	if err = runClientStreamingCall(c, names); err != nil {
		log.Fatalln(err)
	}
	if err = runServerStreamingCall(c, "John", 10); err != nil {
		log.Fatalln(err)
	}
	names = []string{"John", "Paul", "John", "George", "Ringo", "Paul", "John", "Paul", "George", "John"}
	if err = runBidirectionalStreamingCall(c, names); err != nil {
		log.Fatalln(err)
	}
}
