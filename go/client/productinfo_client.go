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
	"google.golang.org/protobuf/types/known/wrapperspb"
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

	// Create new product
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

	// Get new product
	product, err := pic.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	if err != nil {
		log.Fatalf("Could not get product: %v", err)
	}
	log.Printf("Product: %s", product.String())

	// Server-streaming searchOrders
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	searchStream, err := omc.SearchOrders(ctx, &wrappers.StringValue{Value: "tape"})
	if err != nil {
		log.Fatalf("Failed to get searchOrders stream: %v", err)
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

	// Client-streaming updateOrders
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	updateStream, err := omc.UpdateOrders(ctx)
	if err != nil {
		log.Fatalf("Failed to get sendOrders stream: %v", err)
	}

	// Update order 1
	updateOrder1 := pb.Order{Id: proto.String("0001"), Items: []string{"fertilizer", "flower pot", "topsoil"}, Description: proto.String("gardening"), Price: proto.Float32(45.99), Destination: proto.String("Z")}
	if err := updateStream.Send(&updateOrder1); err != nil {
		log.Fatalf("sendOrders send (%v) failed : %v", updateOrder1, err)
	}
	// New order 5
	newOrder5 := pb.Order{Id: proto.String("0005"), Items: []string{"dog food", "rope", "bug zapper"}, Description: proto.String("Sunday regular"), Price: proto.Float32(38.99), Destination: proto.String("Z")}
	if err := updateStream.Send(&newOrder5); err != nil {
		log.Fatalf("sendOrders send (%v) failed : %v", newOrder5, err)
	}
	// Close client side
	updateRes, err := updateStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("sendOrders close failed : %v", err)
	}
	log.Printf("sendOrders response : %v", updateRes)

	// Bidirectional processOrders
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	streamProcOrder, err := omc.ProcessOrders(ctx)
	if err != nil {
		log.Fatalf("Failed to get processOrders stream : %v", err)
	}

	// Send orders to processOrders
	orderIds := []string{"0001", "0002", "0003", "0004"}
	for _, orderId := range orderIds {
		if err := streamProcOrder.Send(&wrapperspb.StringValue{Value: orderId}); err != nil {
			log.Fatalf("processOrders send(%s) failed : %v", orderId, err)
		}
	}

	channel := make(chan struct{})
	// Start up receiver
	go asyncClientBidirectionalRPC(streamProcOrder, channel)
	// Mimic delay
	time.Sleep(time.Millisecond * 1000)
	if err := streamProcOrder.Send(&wrapperspb.StringValue{Value: "0005"}); err != nil {
		log.Fatalf("processOrders send(%s) failed : %v", "0005", err)
	}
	if err := streamProcOrder.CloseSend(); err != nil {
		log.Fatalf("processOrders closeSend failed : %v", err)
	}

	<-channel

	// TODO: Test bidirectional context cancel
}

// asyncClientBidirectionalRPC
func asyncClientBidirectionalRPC(streamProcOrder grpc.BidiStreamingClient[wrapperspb.StringValue, pb.CombinedShipment], c chan struct{}) {
	for {
		combinedShipment, err := streamProcOrder.Recv()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Printf("Failed to receive combined shipment : %v", err)
			}
		}
		log.Print("Combined shipment : ", combinedShipment.GetOrdersList())
	}
	c <- struct{}{}
}
