package grp

import (
	"context"
	"log"
	"net"

	pb "github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/api/grpc/orderpb"
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	service service.OrderService
}

func NewOrderServer(s service.OrderService) *OrderServer {
	return &OrderServer{service: s}
}

func (s *OrderServer) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := s.service.GetOrders()
	if err != nil {
		return nil, err
	}

	var pbOrders []*pb.Order
	for _, o := range orders {
		pbOrders = append(pbOrders, &pb.Order{
			Id:         uint32(o.ID),
			Item:       o.Item,
			Amount:     int32(o.Amount),
			PatientId:  o.PatientID,
			Medication: o.Medication,
			Dosage:     o.Dosage,
			Status:     o.Status,
		})
	}

	return &pb.ListOrdersResponse{Orders: pbOrders}, nil
}

func StartGRPCServer(orderService service.OrderService, port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, NewOrderServer(orderService))
	// habilita reflection
	reflection.Register(grpcServer)

	log.Printf("🚀 gRPC server running on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
