package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var MinRatingToShowDestination = 4.5

var orderDatabase = NewDatabase()
var orderExecuterIndex = NewDatabase()

func handleAssignOrderRequest(w http.ResponseWriter, r *http.Request) {
    orderId := r.URL.Query().Get("order_id")
    executerId := r.URL.Query().Get("executer_id")
    locale := r.URL.Query().Get("locale")

    if orderId == "" || executerId == "" || locale == "" {
        http.Error(w, "Missing parameters, please provide order_id, executer_id, and locale", http.StatusBadRequest)
        return
    }
    
    orderInfo := GetOrderInfo(orderId, executerId)

    order := AssignedOrder{
        AssignOrderId:    GenerateUUID(),
        OrderId:          orderId,
        ExecuterId:       executerId,
        CoinCoeff:        orderInfo.PriceComponents.CoinCoeff,
        CoinBonusAmount:  orderInfo.PriceComponents.BonusAmount,
        FinalCoinAmount:  orderInfo.FinalCoinAmount,
        RouteInformation: "",
        AssignTime:       time.Now(),
        AcquireTime:      time.Time{}, // Assuming this is initially unset
    }

    if orderInfo.ExecutorProfile.Rating >= MinRatingToShowDestination {
        order.RouteInformation = fmt.Sprintf("Order at zone \"%s\"", "TODO_ZONE_NAME")
    } else {
        order.RouteInformation = "Order at somewhere"
    }

    fmt.Printf(">> New order handled! %+v\n", order)

    orderDatabase.AddItem(orderId, order)
    orderExecuterIndex.AddItem(executerId, orderId)
}

func main() {
    http.HandleFunc("/order/assign", handleAssignOrderRequest)

    fmt.Println("Starting server at :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
