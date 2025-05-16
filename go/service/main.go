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
		"0001": {Id: proto.String("0001"), Items: []string{"duct tape", "cow bell"}},
		"0002": {Id: proto.String("0002"), Items: []string{"kitty litter", "tire iron", "cow bell"}},
		"0003": {Id: proto.String("0003"), Items: []string{"nail gun", "thumb tacks", "duct tape"}},
		"0004": {Id: proto.String("0004"), Items: []string{"painter tape", "caulking", "gaffe tape"}},
	}
	pb.RegisterOrderManagementServer(s, &OrderManagementServer{orderMap: orderMap})

	log.Printf("Starting gRPC listener on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
