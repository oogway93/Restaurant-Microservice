package main 

import (
	"context"
	"log"
	"net"

	pb "github.com/oogway93/microservice/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMyServiceServer
}

func (s *server) GetData(ctx context.Context, req *pb.RequestData) (*pb.ResponseData, error) {
	return &pb.ResponseData{
		Id:    req.Id,
		Value: "Data for " + req.Id,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMyServiceServer(s, &server{})
	log.Println("Server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
