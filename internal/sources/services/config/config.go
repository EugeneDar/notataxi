package config

import (
	"context"
	"notataxi/internal/protobufs/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Source struct {
	GRPCConfig config.ConfigServiceClient
}

func NewSource() (*Source, error) {
	configCon, err := grpc.NewClient("mocks-service.wholeservicenamespace.svc.cluster.local:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Source{
		GRPCConfig: config.NewConfigServiceClient(configCon),
	}, err
}

func (s *Source) CallGetConfig(ctx context.Context) (*config.ConfigResponse, error) {
	return s.GRPCConfig.GetConfig(ctx, nil)
}
