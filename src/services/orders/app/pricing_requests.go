package main

// todo(eugenedar): use generated protobuf structs instead of hand-written ones
type OrderInfoResponse struct {
    OrderId         string
    FinalCoinAmount float64
    PriceComponents PriceComponents
    ExecutorProfile ExecutorProfile
}

type ExecutorProfile struct {
    Id    string
    Tags  []string
    Rating float64
}

type PriceComponents struct {
    BaseCoinAmount float64
    CoinCoeff      float64
    BonusAmount    float64
}

func GetOrderInfo(orderId string, executorId string) OrderInfoResponse {
	// todo(eugenedar): implement a real API call here
    return OrderInfoResponse{
		OrderId: orderId,
		FinalCoinAmount: 100.0,
		PriceComponents: PriceComponents{
			BaseCoinAmount: 50.0,
			CoinCoeff:      1.5,
			BonusAmount:    20.0,
		},
		ExecutorProfile: ExecutorProfile{
			Id:     executorId,
			Tags:   []string{"fast", "reliable"},
			Rating: 4.8,
		},
    }
}
