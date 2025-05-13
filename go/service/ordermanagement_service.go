package main

import (
	"fmt"
	"log"
	pb "productinfo/service/ecommerce"
	"strings"

	"github.com/golang/protobuf/ptypes/wrappers"
)

type OrderManagementServer struct {
	orderMap map[string]*pb.Order
	// This is used for forward-compatibilty in the case where the current server does not implement all methods (will return error message with "unimplemented" code)
	pb.UnimplementedOrderManagementServer
}

func (s *OrderManagementServer) searchOrders(searchQuery *wrappers.StringValue, stream pb.OrderManagement_SearchOrdersServer) error {
	for key, order := range s.orderMap {
		log.Print(key, order)
		for _, itemStr := range order.Items {
			log.Print(itemStr)
			if strings.Contains(itemStr, searchQuery.Value) {
				// Send the matching orders in a stream
				err := stream.Send(order)
				if err != nil {
					return fmt.Errorf("error sending message to stream : %v", err)
				}
				log.Printf("Matching Order found : " + key)
				break
			}
		}
	}
	return nil
}
