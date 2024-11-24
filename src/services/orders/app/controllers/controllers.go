package controllers

import (
	"app/src/services/orders/app/database"
	"app/src/services/sources/protobufs/sources"
	"context"

	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	"github.com/gin-gonic/gin"
)

type Service struct {
	GRPCSourcesClient sources.SourcesServiceClient
}

func NewService() (*Service, error) {
	conSources, err := grpc.Dial("sources:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	return &Service{
		GRPCSourcesClient: sources.NewSourcesServiceClient(conSources),
	}, nil
}

func (s *Service) AssignOrderRequestHandler(c *gin.Context) {
	orderId := c.Query("order_id")
	executorId := c.Query("executor_id")
	zoneId := c.Query("zone_id")
	if orderId == "" || executorId == "" || zoneId == "" {
		c.JSON(400, gin.H{"message": "Missing parameters, please provide order_id, executor_id and zone_id"})
		return
	}

	ctx := context.Background()
	orderProfile, err := s.GRPCSourcesClient.GetOrderInfo(ctx, &sources.SourcesRequest{
		OrderId:    orderId,
		ExecutorId: executorId,
	})

	if err != nil {
		log.Println(err)
		return
	}

	if err = database.AddAssignedOrder(orderProfile); err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"assigned_orders_orderid_key\"" {
			c.JSON(400, gin.H{"message": "AssignedOrder with provided orderId already exists"})
			log.Printf("AssignedOrder with provided orderId=%s already exists\n", orderId)
			return
		}
		c.JSON(500, gin.H{"message": "Unknown error"})
		log.Printf("Error executing AddAssignedOrder: %s\n", err.Error())
		return
	}

	c.JSON(200, gin.H{"message": "Successfully created"})
	log.Printf("[Info] New order have been handled! Order profile:\n%+v\n", orderProfile)
}

func (s *Service) AcquireOrderRequestHandler(c *gin.Context) {
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

func (s *Service) CancelOrderRequestHandler(c *gin.Context) {
	orderId := c.Query("order_id")
	if orderId == "" {
		c.JSON(400, gin.H{"message": "Missing parameters, please provide order_id"})
		return
	}

	found, err := database.CancelAssignedOrder(orderId)
	if err != nil {
		c.JSON(500, gin.H{"message": "Unknown error"})
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

func (s *Service) CleanDatabaseRequestHandler(c *gin.Context) {
	if err := database.CleanDatabase(); err != nil {
		c.JSON(500, gin.H{"error": "Failed to clean database"})
		log.Printf("Error executing CleanDatabase: %s\n", err.Error())
		return
	}
	c.JSON(200, gin.H{"message": "Database cleaned successfully"})
}

func (s *Service) CleanTestOrdersHandler(c *gin.Context) {
	if err := database.CleanTestOrders(); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to clean test orders: %v", err)})
		log.Printf("Error executing CleanTestOrders: %s\n", err.Error())
		return
	}
	c.JSON(200, gin.H{"message": "Test orders cleaned successfully"})
}
