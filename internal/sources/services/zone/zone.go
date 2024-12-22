package zone

import (
	"context"
	"log"
	"notataxi/internal/protobufs/zone_data"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Source struct {
	GRPCZone zone_data.ZoneDataServiceClient
}

func NewSource() (*Source, error) {
	zoneCon, err := grpc.NewClient("mocks-service.wholeservicenamespace.svc.cluster.local:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	return &Source{
		GRPCZone: zone_data.NewZoneDataServiceClient(zoneCon),
	}, err
}

func (s *Source) GetZoneData(ctx context.Context, zoneId string) (*zone_data.ZoneDataResponse, error) {
	return s.GRPCZone.GetZoneData(ctx, &zone_data.ZoneDataRequest{ZoneId: zoneId})
}
