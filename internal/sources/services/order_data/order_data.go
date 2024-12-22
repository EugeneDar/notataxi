package order_data

import (
	"notataxi/internal/protobufs/order_data"

	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Source struct {
	GRPCOrderData order_data.OrderDataServiceClient
}

func NewSource() (*Source, error) {
	orderDataCon, err := grpc.NewClient("mocks-service.wholeservicenamespace.svc.cluster.local:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	return &Source{
		GRPCOrderData: order_data.NewOrderDataServiceClient(orderDataCon),
	}, err
}

func (s *Source) GetOrderData(ctx context.Context, orderId string) (*order_data.OrderDataResponse, error) {
	return s.GRPCOrderData.GetOrderData(ctx, &order_data.OrderDataRequest{OrderId: orderId})
}
