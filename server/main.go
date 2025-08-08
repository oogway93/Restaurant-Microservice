package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/oogway93/microservice/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedRestaurantServiceServer
	db *mongo.Collection
}

type Order struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Title string             `bson:"title"`
	Price int64              `bson:"price"`
}

func (s *server) CreateOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
	res, err := s.db.InsertOne(ctx, Order{
		Title: req.GetTitle(),
		Price: req.GetPrice(),
	})
	if err != nil {
		return nil, err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("invalid object id")
	}
	return &pb.OrderResponse{
		Id:    oid.Hex(),
		Title: req.GetTitle(),
		Price: req.GetPrice(),
	}, nil
}

func (s *server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}

	// Query MongoDB
	var order Order
	err = s.db.FindOne(ctx, bson.M{"_id": oid}).Decode(&order)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	return &pb.OrderResponse{
		Id:    order.ID.Hex(),
		Title: order.Title,
		Price: order.Price,
	}, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@mongo:27017"))
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}
	defer client.Disconnect(ctx)
	
	db := client.Database("mydb").Collection("orders")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRestaurantServiceServer(s, &server{db: db})
	log.Println("Server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
