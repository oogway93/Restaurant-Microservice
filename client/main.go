package main 

import (
	"context"
	"log"
	"time"

	pb "github.com/oogway93/microservice/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewMyServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetData(ctx, &pb.RequestData{Id: "123"})
	if err != nil {
		log.Fatalf("could not get data: %v", err)
	}
	log.Printf("Response: %s", r.GetValue())
}
