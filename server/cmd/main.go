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

func (s *CallServer) Call(ctx context.Context, in *pb.CallRequest) (*pb.CallResponse, error) {
	resp := &pb.CallResponse{}
	resp.Message = fmt.Sprintf("Hello. I'm %s.", in.GetName())
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
