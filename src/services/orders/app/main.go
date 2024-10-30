package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
    "app/database"
    "app/requests"
    "app/model"
    "app/utils"
)

var MinRatingToShowDestination = 4.5

var orderDatabase = database.NewDatabase()
var orderExecuterIndex = database.NewDatabase()

func handleAssignOrderRequest(w http.ResponseWriter, r *http.Request) {
    orderId := r.URL.Query().Get("order_id")
    executerId := r.URL.Query().Get("executer_id")
    locale := r.URL.Query().Get("locale")

    if orderId == "" || executerId == "" || locale == "" {
        http.Error(w, "Missing parameters, please provide order_id, executer_id, and locale", http.StatusBadRequest)
        return
    }
    
    orderInfo := requests.GetOrderInfo(orderId, executerId)

    order := model.AssignedOrder{
        AssignOrderId:    utils.GenerateUUID(),
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

func handleAcquireOrderRequest(w http.ResponseWriter, r *http.Request) {
    executerId := r.URL.Query().Get("executer_id")

    if executerId == "" {
        http.Error(w, "Missing parameters, please provide executer_id", http.StatusBadRequest)
        return
    }

    orderId, err := orderExecuterIndex.GetItem(executerId)
    if err != nil {
        fmt.Printf("Order for executer ID \"%s\" not found!\n", executerId)
        return
    }

    orderData, err := orderDatabase.GetItem(orderId.(string))
    if err != nil {
        fmt.Printf("Order data for order ID \"%s\" not found!\n", orderId)
        return
    }

    order := orderData.(model.AssignedOrder)
    order.AcquireTime = time.Now()

    fmt.Printf(">> Order acquired! Acquire time == %v\n", order.AcquireTime.Sub(order.AssignTime))
    fmt.Printf(">> Order data: %+v\n", order)
}

func handleCancelOrderRequest(w http.ResponseWriter, r *http.Request) {
    orderId := r.URL.Query().Get("order_id")

    if orderId == "" {
        http.Error(w, "Missing parameters, please provide order_id", http.StatusBadRequest)
        return
    }

    orderData, err := orderDatabase.DeleteItem(orderId)
    if err != nil {
        fmt.Printf("Order data for order ID \"%s\" not found!\n", orderId)
        return
    }

    order := orderData.(model.AssignedOrder)
    orderExecuterIndex.DeleteItem(order.ExecuterId)

    fmt.Printf(">> Order was cancelled! %+v\n", order)
}

func main() {
    http.HandleFunc("/order/assign", handleAssignOrderRequest)
    http.HandleFunc("/order/acquire", handleAcquireOrderRequest)
    http.HandleFunc("/order/cancel", handleCancelOrderRequest);

    fmt.Println("Starting server at :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
