package toll_roads

import (
	"notataxi/internal/protobufs/toll_roads"

	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Source struct {
	GRPCTollRoads toll_roads.TollRoadsServiceClient
}

func NewSource() (*Source, error) {
	tollRoadsCon, err := grpc.NewClient("mocks-service.wholeservicenamespace.svc.cluster.local:9093", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	return &Source{
		GRPCTollRoads: toll_roads.NewTollRoadsServiceClient(tollRoadsCon),
	}, err
}

func (s *Source) GetTollRoads(ctx context.Context, zoneName string) (*toll_roads.TollRoadsResponse, error) {
	return s.GRPCTollRoads.GetTollRoads(ctx, &toll_roads.TollRoadsRequest{DisplayName: zoneName})
}
