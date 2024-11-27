package requests

import (
	"app/src/services/orders/model"
	"app/src/services/sources/protobufs/sources"
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var sourcesService sources.SourcesServiceClient

func ConnectionToSourcesService() (err error) {
	conSources, err := grpc.NewClient("sources-service.wholeservicenamespace.svc.cluster.local:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	sourcesService = sources.NewSourcesServiceClient(conSources)
	return nil
}

func GetOrderInfo(orderId string, executorId string, zoneId string) (*model.AssignedOrder, error) {
	ctx := context.Background()
	response, err := sourcesService.GetOrderInfo(ctx, &sources.SourcesRequest{
		OrderId:    orderId,
		ExecutorId: executorId,
		// ZoneId: zoneId,  // TODO: Why sources service does not require zoneId but provides ZoneName?
	})
	if err != nil {
		return nil, err
	}
	return &model.AssignedOrder{
		AssignedOrderId:             uuid.New().String(),
		OrderId:                     orderId,
		ExecutorId:                  executorId,
		ExecutionStatus:             "assigned",
		CoinCoefficient:             response.GetPriceComponents().GetCoinCoeff(),
		CoinBonusAmount:             response.GetPriceComponents().GetBonusAmount(),
		FinalCoinAmount:             response.GetFinalCoinAmount(),
		ZoneName:                    response.GetZoneDisplayName(),
		HasExecutorFallbackBeenUsed: response.GetUsedExecutorFallback(),
	}, nil
}
