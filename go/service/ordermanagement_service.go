package main

import (
	"fmt"
	"io"
	"log"
	pb "productinfo/service/ecommerce"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type OrderManagementServer struct {
	orderMap       map[string]*pb.Order
	orderBatchSize uint
	shipmentId     uint
	// This is used for forward-compatibilty in the case where the current server does not implement all methods (will return error message with "unimplemented" code)
	pb.UnimplementedOrderManagementServer
}

func (s *OrderManagementServer) SearchOrders(searchQuery *wrapperspb.StringValue, stream grpc.ServerStreamingServer[pb.Order]) error {
	for key, order := range s.orderMap {
		log.Print(key, " ", order)
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

func (s *OrderManagementServer) UpdateOrders(stream grpc.ClientStreamingServer[pb.Order, wrapperspb.StringValue]) error {
	ordersStr := "Updated Order IDs : "
	for {
		order, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(
					&wrapperspb.StringValue{Value: "Orders processed " + ordersStr})
			}
		}

		s.orderMap[order.GetId()] = order

		log.Print("Order ID ", order.GetId(), ": Updated")
		ordersStr += *order.Id + ", "
	}
}

/*
ProcessOrders receives orders from the client, packages them into groups based on shipping destination, and then
sends the groups once the batch size is reached, or the client stops sending.
*/
func (s *OrderManagementServer) ProcessOrders(stream grpc.BidiStreamingServer[wrapperspb.StringValue, pb.CombinedShipment]) error {
	var batchMarker uint = 1
	combinedShipmentMap := make(map[string]*pb.CombinedShipment)
	for {
		orderId, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				for _, comb := range combinedShipmentMap {
					stream.Send(comb)
				}
				return nil
			} else {
				return err
			}
		}

		// Find requested order
		o, exists := s.orderMap[orderId.GetValue()]
		if !exists {
			continue
		}

		// Get shipment if it exists
		shipment, exists := combinedShipmentMap[o.GetDestination()]
		if !exists {
			shipment = &pb.CombinedShipment{Id: proto.String(fmt.Sprintf("%4d", s.shipmentId)), Status: proto.String("OK"), OrdersList: make([]*pb.Order, s.orderBatchSize/2)}
			s.shipmentId++
			combinedShipmentMap[o.GetDestination()] = shipment
		}

		shipment.OrdersList = append(shipment.OrdersList, o)

		if batchMarker == s.orderBatchSize {
			for _, comb := range combinedShipmentMap {
				stream.Send(comb)
			}
			clear(combinedShipmentMap)
			batchMarker = 1
		} else {
			batchMarker++
		}
	}
}
