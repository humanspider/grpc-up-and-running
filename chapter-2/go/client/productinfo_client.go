package main

import (
	"context"
	"log"
	"time"

	pb "productinfo/client/ecommerce"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	name := "Apple iPhone 11"
	description := "Meet Apple iPhone 11. All-new dual camera system with Ultra Wide and Night mode."
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.AddProduct(ctx,
		&pb.Product{Name: proto.String(name), Description: proto.String(description), Price: proto.Float32(price)})
	if err != nil {
		log.Fatalf("Could not add product: %v", err)
	}
	log.Printf("Product ID: %s added successfully", r.Value)

}
