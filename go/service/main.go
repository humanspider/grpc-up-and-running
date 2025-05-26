package main

import (
	"log"
	"net"

	pb "productinfo/service/ecommerce"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &ProductInfoServer{})
	orderMap := map[string]*pb.Order{
		"0001": {Id: proto.String("0001"), Items: []string{"duct tape", "cow bell"}, Destination: proto.String("X")},
		"0002": {Id: proto.String("0002"), Items: []string{"kitty litter", "tire iron", "cow bell"}, Destination: proto.String("X")},
		"0003": {Id: proto.String("0003"), Items: []string{"nail gun", "thumb tacks", "duct tape"}, Destination: proto.String("Y")},
		"0004": {Id: proto.String("0004"), Items: []string{"painter tape", "caulking", "gaffe tape"}, Destination: proto.String("Z")},
	}
	pb.RegisterOrderManagementServer(s, &OrderManagementServer{orderMap: orderMap, orderBatchSize: 3})

	log.Printf("Starting gRPC listener on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
