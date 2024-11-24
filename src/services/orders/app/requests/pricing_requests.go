package requests

import (
	"app/src/services/orders/app/model"
	"app/src/services/orders/app/utils"
)

func GetOrderInfo(orderId string, executorId string, zoneId string) model.AssignedOrder {
	// TODO: implement the gRPC call to the Price service
	return model.AssignedOrder{
		AssignedOrderId:             utils.GenerateUUID(),
		OrderId:                     orderId,
		ExecutorId:                  executorId,
		ExecutionStatus:             "assigned",
		CoinCoefficient:             1.5,
		CoinBonusAmount:             20.0,
		FinalCoinAmount:             50.0,
		ZoneName:                    "Lyubertsy",
		HasExecutorFallbackBeenUsed: false,
	}
}
