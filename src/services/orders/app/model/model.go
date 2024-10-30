package model

import "time"

// todo(eugenedar): use generated protobuf structs instead of hand-written ones
type AssignedOrder struct {
	AssignOrderId     string    `json:"assign_order_id"`
	OrderId           string    `json:"order_id"`
	ExecuterId        string    `json:"executer_id"`
	CoinCoeff         float64   `json:"coin_coeff"`
	CoinBonusAmount   float64   `json:"coin_bonus_amount"`
	FinalCoinAmount   float64   `json:"final_coin_amount"`
	RouteInformation  string    `json:"route_information"`
	AssignTime        time.Time `json:"assign_time"`
	AcquireTime       time.Time `json:"acquire_time"`
}


