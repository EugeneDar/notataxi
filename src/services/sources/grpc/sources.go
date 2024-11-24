package grpcsources

import (
	"app/src/services/sources/protobufs/config"
	"app/src/services/sources/protobufs/executor_profile"
	"app/src/services/sources/protobufs/order_data"
	"app/src/services/sources/protobufs/sources"
	"app/src/services/sources/protobufs/toll_roads"
	"app/src/services/sources/protobufs/zone_data"
	"log"

	"google.golang.org/grpc/credentials/insecure"

	"context"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"google.golang.org/grpc"
)

const (
	ConfigFieldsCount    = 1
	TollRoadsFieldsCount = 1
	ZoneFieldsCount      = 1

	ConfigCacheTTL = time.Minute * 1
	BaseCacheTTL   = time.Minute * 10

	ExecutorTimeOut = time.Second * 5
)

type ServiceAPI struct {
	sources.UnsafeSourcesServiceServer

	GRPCConfig config.ConfigServiceClient

	GRPCOrderData order_data.OrderDataServiceClient
	GRPCZone      zone_data.ZoneDataServiceClient
	GRPCTollRoads toll_roads.TollRoadsServiceClient

	GRPCExecutor         executor_profile.ExecutorProfileServiceClient
	GRPCExecutorFallback executor_profile.ExecutorProfileServiceClient

	ConfigCache    *expirable.LRU[string, *config.ConfigResponse]
	TollRoadsCache *expirable.LRU[string, *toll_roads.TollRoadsResponse]
	ZoneCache      *expirable.LRU[string, *zone_data.ZoneDataResponse]
}

func Register(gRPC *grpc.Server) error {
	configCon, err := grpc.NewClient("mocks-service.default.svc.cluster.local:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	orderDataCon, err := grpc.NewClient("mocks-service.default.svc.cluster.local:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	zoneCon, err := grpc.NewClient("mocks-service.default.svc.cluster.local:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	tollRoadsCon, err := grpc.NewClient("mocks-service.default.svc.cluster.local:9093", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	executorCon, err := grpc.NewClient("mocks-service.default.svc.cluster.local:9094", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	executorFallbackCon, err := grpc.NewClient("localhost:9095", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	sources.RegisterSourcesServiceServer(gRPC, &ServiceAPI{
		GRPCConfig: config.NewConfigServiceClient(configCon),

		GRPCOrderData: order_data.NewOrderDataServiceClient(orderDataCon),
		GRPCZone:      zone_data.NewZoneDataServiceClient(zoneCon),
		GRPCTollRoads: toll_roads.NewTollRoadsServiceClient(tollRoadsCon),

		GRPCExecutor:         executor_profile.NewExecutorProfileServiceClient(executorCon),
		GRPCExecutorFallback: executor_profile.NewExecutorProfileServiceClient(executorFallbackCon),

		ConfigCache:    expirable.NewLRU[string, *config.ConfigResponse](ConfigFieldsCount, nil, ConfigCacheTTL),
		TollRoadsCache: expirable.NewLRU[string, *toll_roads.TollRoadsResponse](TollRoadsFieldsCount, nil, BaseCacheTTL),
		ZoneCache:      expirable.NewLRU[string, *zone_data.ZoneDataResponse](ZoneFieldsCount, nil, BaseCacheTTL),
	})

	return nil
}

func (s *ServiceAPI) PriceCalculate(ctx context.Context, baseCoinAmount, bonusAmount int32, coinCoeff float32) (int32, error) {
	configInfo, ok := s.ConfigCache.Get("config")
	if !ok {
		configInfo, err := s.GRPCConfig.GetConfig(ctx, nil)
		if err != nil {
			return 0, err
		}
		s.ConfigCache.Add("config", configInfo)
	}
	return max(int32(float32(baseCoinAmount)*coinCoeff)+bonusAmount, configInfo.GetMinPrice()), nil
}

func (s *ServiceAPI) GetOrderInfo(ctx context.Context, req *sources.SourcesRequest) (*sources.SourcesResponse, error) {
	orderDataInfo, err := s.GRPCOrderData.GetOrderData(ctx, &order_data.OrderDataRequest{OrderId: req.GetOrderId()})
	if err != nil {
		return nil, err
	}

	zoneInfo, ok := s.ZoneCache.Get(orderDataInfo.GetZoneId())
	if !ok {
		zoneInfo, err = s.GRPCZone.GetZoneData(ctx, &zone_data.ZoneDataRequest{ZoneId: orderDataInfo.GetZoneId()})
		if err != nil {
			return nil, err
		}
		s.ZoneCache.Add(orderDataInfo.GetZoneId(), zoneInfo)
	}

	tollRoadsInfo, ok := s.TollRoadsCache.Get(zoneInfo.GetDisplayName())
	if !ok {
		tollRoadsInfo, err = s.GRPCTollRoads.GetTollRoads(ctx, &toll_roads.TollRoadsRequest{DisplayName: zoneInfo.GetDisplayName()})
		if err != nil {
			return nil, err
		}
		s.TollRoadsCache.Add(zoneInfo.GetDisplayName(), tollRoadsInfo)
	}

	executorContext, cancel := context.WithTimeout(ctx, ExecutorTimeOut)
	defer cancel()

	useFallback := false

	executorInfo, err := s.GRPCExecutor.GetExecutorProfile(executorContext, &executor_profile.ExecutorProfileRequest{DisplayName: zoneInfo.GetDisplayName()})
	if err != nil {
		useFallback = true
		executorInfo, err = s.GRPCExecutorFallback.GetExecutorProfile(ctx, &executor_profile.ExecutorProfileRequest{DisplayName: zoneInfo.GetDisplayName()})
		if err != nil {
			return nil, err
		}
	}

	finalPrice, err := s.PriceCalculate(ctx, orderDataInfo.GetBaseCoinAmount(), tollRoadsInfo.GetBonusAmount(), zoneInfo.GetCoinCoeff())
	if err != nil {
		return nil, err
	}

	resp := sources.SourcesResponse{
		OrderId:         req.GetOrderId(),
		FinalCoinAmount: finalPrice,
		PriceComponents: &sources.PriceComponents{
			BaseCoinAmount: orderDataInfo.GetBaseCoinAmount(),
			CoinCoeff:      zoneInfo.GetCoinCoeff(),
			BonusAmount:    tollRoadsInfo.GetBonusAmount(),
		},
		ExecutorProfile: &sources.ExecutorProfile{
			Id:     executorInfo.GetId(),
			Tags:   executorInfo.GetTags(),
			Rating: executorInfo.GetRating(),
		},
		ZoneDisplayName:      zoneInfo.GetDisplayName(),
		UsedExecutorFallback: useFallback,
	}

	return &resp, nil
}
