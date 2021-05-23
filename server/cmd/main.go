package main

import (
	"errors"
	"fmt"
	"grpc-lesson/gen/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

const port=":50001"

type Server struct {
	pb.UnimplementedCallMeJohnServer
}

func (s *Server) Call(req *pb.CallRequest) (*pb.CallResponse, error) {
	resp := &pb.CallResponse{}
	switch req.GetName() {
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
	var server Server
	pb.RegisterCallMeJohnServer(s, &server)
	if err := s.Serve(lis); err != nil {
		return errors.New("serve is failed")
	}
	return nil
}
