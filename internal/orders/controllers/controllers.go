package controllers

import (
	"notataxi/internal/orders/database"
	"notataxi/internal/orders/requests"

	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func AssignOrderRequestHandler(c *gin.Context) {
	orderId := c.Query("order_id")
	executorId := c.Query("executor_id")
	zoneId := c.Query("zone_id")
	if orderId == "" || executorId == "" || zoneId == "" {
		c.JSON(400, gin.H{"message": "Missing parameters, please provide order_id, executor_id and zone_id"})
		return
	}

	orderProfile, err := requests.GetOrderInfo(orderId, executorId, zoneId)
	if err != nil {
		c.JSON(500, gin.H{"message": "Unknown error"})
		log.Printf("Error requesting sources service: %s\n", err.Error())
		return
	}

	created, err := database.AddAssignedOrder(orderProfile)
	if err != nil {
		c.JSON(500, gin.H{"message": "Unknown error"})
		log.Printf("Error executing AddAssignedOrder: '%s'\norderProfile: %v", err.Error(), orderProfile)
		return
	}
	if !created {
		c.JSON(400, gin.H{"message": "AssignedOrder with provided orderId already exists"})
		log.Printf("AssignedOrder with provided orderId=%s already exists\n", orderId)
		return
	}

	c.JSON(200, gin.H{"message": "Successfully created"})
	log.Printf("[Info] New order have been handled! Order profile:\n%+v\n", orderProfile)
}

func AcquireOrderRequestHandler(c *gin.Context) {
	executorId := c.Query("executor_id")
	if executorId == "" {
		c.JSON(400, gin.H{"message": "Missing parameters, please provide executor_id"})
		return
	}

	orderProfile, err := database.AcquireAssignedOrder(executorId)
	if err != nil {
		c.JSON(500, gin.H{"message": "Unknown error"})
		log.Printf("Error executing AcquireAssignedOrder: %s\n", err.Error())
		return
	}
	if orderProfile == nil {
		c.JSON(200, gin.H{"message": "There are no orders assigned to you"})
		log.Printf("[Info] There are no orders assigned to ExecutorId=%s\n", executorId)
		return
	}

	c.JSON(200, gin.H{"message": "Successfully acquired", "order_profile": orderProfile})
	log.Printf("[Info] Order has just been acquired! Order profile:\n%+v\n", orderProfile)
}

func CancelOrderRequestHandler(c *gin.Context) {
	orderId := c.Query("order_id")
	if orderId == "" {
		c.JSON(400, gin.H{"message": "Missing parameters, please provide order_id"})
		return
	}

	found, err := database.CancelAssignedOrder(orderId)
	if err != nil {
		c.JSON(500, gin.H{})
		log.Printf("Error executing CancelAssignedOrder: %s\n", err.Error())
		return
	}
	if !found {
		c.JSON(200, gin.H{"message": fmt.Sprintf("AssignedOrder with OrderId %s does not exist or has already been canceled", orderId)})
		log.Printf("[Info] AssignedOrder with OrderId %s does not exist or has already been canceled or completed\n", orderId)
		return
	}

	c.JSON(200, gin.H{"message": "Successfully canceled"})
	log.Printf("[Info] Have just cancel order with OrderId=%s\n", orderId)
}

func CleanDatabaseRequestHandler(c *gin.Context) {
	if err := database.CleanDatabase(); err != nil {
		c.JSON(500, gin.H{"error": "Failed to clean database"})
		log.Printf("Error executing CleanDatabase: %s\n", err.Error())
		return
	}
	c.JSON(200, gin.H{"message": "Database cleaned successfully"})
}

func CleanTestOrdersHandler(c *gin.Context) {
	if err := database.CleanTestOrders(); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to clean test orders: %v", err)})
		log.Printf("Error executing CleanTestOrders: %s\n", err.Error())
		return
	}
	c.JSON(200, gin.H{"message": "Test orders cleaned successfully"})
}
