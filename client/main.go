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

	c := pb.NewRestaurantServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.CreateOrder(ctx, &pb.OrderRequest{
		Title:  "GPU 5070 Super",
		Price: 100000,
	})
	if err != nil {
		log.Fatalf("CreateOrder failed: %v", err)
	}
	log.Printf("Created Order: ID=%s, Title=%s", res.Id, res.Title)

	// Get User
	user, err := c.GetOrder(ctx, &pb.GetOrderRequest{Id: res.Id})
	if err != nil {
		log.Fatalf("GetOrder failed: %v", err)
	}
	log.Printf("Retrieved Order: %v", user)
}
