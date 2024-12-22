package executor

import (
	"notataxi/internal/protobufs/executor_profile"

	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Source struct {
	GRPCExecutor         executor_profile.ExecutorProfileServiceClient
	GRPCExecutorFallback executor_profile.ExecutorProfileServiceClient
}

func NewSource() (*Source, error) {
	executorCon, err := grpc.NewClient("mocks-service.wholeservicenamespace.svc.cluster.local:9094", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	executorFallbackCon, err := grpc.NewClient("mocks-service.wholeservicenamespace.svc.cluster.local:9095", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	return &Source{
		GRPCExecutor:         executor_profile.NewExecutorProfileServiceClient(executorCon),
		GRPCExecutorFallback: executor_profile.NewExecutorProfileServiceClient(executorFallbackCon),
	}, err
}

func (s *Source) GetExecutorProfile(ctx context.Context, executorId string) (*executor_profile.ExecutorProfileResponse, error) {
	return s.GRPCExecutor.GetExecutorProfile(ctx, &executor_profile.ExecutorProfileRequest{ExecutorId: executorId})
}

func (s *Source) GetExecutorProfileWithFallback(ctx context.Context, executorId string) (*executor_profile.ExecutorProfileResponse, error, bool) {
	executorInfo, err := s.GRPCExecutor.GetExecutorProfile(ctx, &executor_profile.ExecutorProfileRequest{ExecutorId: executorId})
	if err != nil {
		executorInfo, err = s.GRPCExecutorFallback.GetExecutorProfile(ctx, &executor_profile.ExecutorProfileRequest{ExecutorId: executorId})
		if err != nil {
			return nil, err, true
		}
		return executorInfo, nil, true
	}
	return executorInfo, nil, false
}
