package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "productinfo/client/ecommerce"

	"github.com/golang/protobuf/ptypes/wrappers"
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
	pic := pb.NewProductInfoClient(conn)
	omc := pb.NewOrderManagementClient(conn)

	name := "Apple iPhone 11"
	description := "Meet Apple iPhone 11. All-new dual camera system with Ultra Wide and Night mode."
	price := float32(1000.0)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := pic.AddProduct(ctx,
		&pb.Product{Name: proto.String(name), Description: proto.String(description), Price: proto.Float32(price)})
	if err != nil {
		log.Fatalf("Could not add product: %v", err)
	}
	log.Printf("Product ID: %s added successfully", *r.Value)

	product, err := pic.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	if err != nil {
		log.Fatalf("Could not get product: %v", err)
	}
	log.Printf("Product: %s", product.String())

	ctx, cancel = context.WithCancel(context.Background())
	searchStream, err := omc.SearchOrders(ctx, &wrappers.StringValue{Value: "tape"})
	if err != nil {
		log.Fatalf("Could not get searchOrders stream: %v", err)
	}

	for {
		searchOrder, err := searchStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatalf("Failed to receive order: %v", err)
			}
		}
		log.Print("Search Result : ", searchOrder)
	}
}
