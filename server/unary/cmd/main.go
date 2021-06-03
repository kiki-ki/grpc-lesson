package main

import (
	"context"
	"errors"
	"fmt"
	"grpc-lesson/gen/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

const port = ":50051"

type CallServer struct {
	pb.UnimplementedCallServer
}

func (s *CallServer) CallMeJohn(ctx context.Context, in *pb.CallMeJohnRequest) (*pb.CallMeJohnResponse, error) {
	resp := &pb.CallMeJohnResponse{}
	switch in.GetName() {
	case "john", "JOHN", "John":
		resp.Message = "Hi."
	default:
		resp.Message = "Please call me John."
	}
	return resp, nil
}

func main() {
	fmt.Printf("server is listening on port%s\n", port)
	if err := set(); err != nil {
		log.Fatalln(err.Error())
	}
}

func set() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}
	s := grpc.NewServer()
	pb.RegisterCallServer(s, &CallServer{})
	if err := s.Serve(lis); err != nil {
		return errors.New("serve is failed")
	}
	return nil
}
