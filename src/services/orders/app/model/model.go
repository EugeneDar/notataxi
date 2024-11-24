package model

import "time"

type AssignedOrder struct {
	AssignedOrderId             string    `json:"assigned_order_id"`
	OrderId                     string    `json:"order_id"`
	ExecutorId                  string    `json:"executor_id"`
	ExecutionStatus             string    `json:"execution_status"`
	CoinCoefficient             float64   `json:"coin_coefficient"`
	CoinBonusAmount             float64   `json:"coin_bonus_amount"`
	FinalCoinAmount             float64   `json:"final_coin_amount"`
	ZoneName                    string    `json:"zone_name"`
	HasExecutorFallbackBeenUsed bool      `json:"has_executor_fallback_been_used"`
	AssignTime                  time.Time `json:"assign_time"`
	LastAcquireTime             time.Time `json:"last_acquire_time"`
}
