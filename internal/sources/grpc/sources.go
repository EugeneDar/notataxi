package grpcsources

import (
	grpc_config "notataxi/internal/protobufs/config"
	"notataxi/internal/protobufs/executor_profile"
	"notataxi/internal/protobufs/sources"
	grpc_toll_roads "notataxi/internal/protobufs/toll_roads"
	"notataxi/internal/protobufs/zone_data"
	"notataxi/internal/sources/services/config"
	"notataxi/internal/sources/services/executor"
	"notataxi/internal/sources/services/order_data"
	"notataxi/internal/sources/services/toll_roads"
	"notataxi/internal/sources/services/zone"

	"context"
	"log"
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

	Config *config.Source

	OrderData *order_data.Source
	Zone      *zone.Source
	TollRoads *toll_roads.Source

	Executor *executor.Source

	ConfigCache    *expirable.LRU[string, *grpc_config.ConfigResponse]
	TollRoadsCache *expirable.LRU[string, *grpc_toll_roads.TollRoadsResponse]
	ZoneCache      *expirable.LRU[string, *zone_data.ZoneDataResponse]
}

func Register(gRPC *grpc.Server) error {
	configSource, err := config.NewSource()
	if err != nil {
		return err
	}

	orderDataSource, err := order_data.NewSource()
	if err != nil {
		log.Fatal(err)
	}

	zoneSource, err := zone.NewSource()
	if err != nil {
		log.Fatal(err)
	}

	tollRoadsSource, err := toll_roads.NewSource()
	if err != nil {
		log.Fatal(err)
	}

	executorSource, err := executor.NewSource()
	if err != nil {
		log.Fatal(err)
	}

	sources.RegisterSourcesServiceServer(gRPC, &ServiceAPI{
		Config: configSource,

		OrderData: orderDataSource,
		Zone:      zoneSource,
		TollRoads: tollRoadsSource,

		Executor: executorSource,

		ConfigCache:    expirable.NewLRU[string, *grpc_config.ConfigResponse](ConfigFieldsCount, nil, ConfigCacheTTL),
		TollRoadsCache: expirable.NewLRU[string, *grpc_toll_roads.TollRoadsResponse](TollRoadsFieldsCount, nil, BaseCacheTTL),
		ZoneCache:      expirable.NewLRU[string, *zone_data.ZoneDataResponse](ZoneFieldsCount, nil, BaseCacheTTL),
	})

	return nil
}

func (s *ServiceAPI) PriceCalculate(ctx context.Context, baseCoinAmount, bonusAmount int32, coinCoeff float32) (int32, error) {
	configInfo, ok := s.ConfigCache.Get("config")
	if !ok {
		configInfo, err := s.Config.CallGetConfig(ctx)
		if err != nil {
			return 0, err
		}
		s.ConfigCache.Add("config", configInfo)
	}
	return max(int32(float32(baseCoinAmount)*coinCoeff)+bonusAmount, configInfo.GetMinPrice()), nil
}

func (s *ServiceAPI) GetOrderInfo(ctx context.Context, req *sources.SourcesRequest) (*sources.SourcesResponse, error) {
	executorErrCh := make(chan error)
	executorInfoCh := make(chan *executor_profile.ExecutorProfileResponse)
	useFallbackCh := make(chan bool)

	executorContext, cancel := context.WithTimeout(ctx, ExecutorTimeOut)
	defer cancel()

	go func() {
		res, err, fallback := s.Executor.GetExecutorProfileWithFallback(executorContext, req.GetExecutorId())
		executorInfoCh <- res
		executorErrCh <- err
		useFallbackCh <- fallback
	}()

	orderDataInfo, err := s.OrderData.GetOrderData(ctx, req.GetOrderId())
	if err != nil {
		return nil, err
	}

	zoneInfo, ok := s.ZoneCache.Get(orderDataInfo.GetZoneId())
	if !ok {
		zoneInfo, err = s.Zone.GetZoneData(ctx, orderDataInfo.GetZoneId())
		if err != nil {
			return nil, err
		}
		s.ZoneCache.Add(orderDataInfo.GetZoneId(), zoneInfo)
	}

	tollRoadsInfo, ok := s.TollRoadsCache.Get(zoneInfo.GetDisplayName())
	if !ok {
		tollRoadsInfo, err = s.TollRoads.GetTollRoads(ctx, zoneInfo.GetDisplayName())
		if err != nil {
			return nil, err
		}
		s.TollRoadsCache.Add(zoneInfo.GetDisplayName(), tollRoadsInfo)
	}

	finalPrice, err := s.PriceCalculate(ctx, orderDataInfo.GetBaseCoinAmount(), tollRoadsInfo.GetBonusAmount(), zoneInfo.GetCoinCoeff())
	if err != nil {
		return nil, err
	}

	if executorErr := <-executorErrCh; executorErr != nil {
		return nil, executorErr
	}
	executorInfo := <-executorInfoCh
	useFallback := <-useFallbackCh

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
