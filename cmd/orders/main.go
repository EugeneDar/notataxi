package main

import (
	"notataxi/internal/orders/controllers"
	"notataxi/internal/orders/database"
	"notataxi/internal/orders/requests"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	err := database.EstablishConnection()
	for err != nil {
		log.Printf("Error connecting to database: %s; retrying...\n", err.Error())
		err = database.EstablishConnection()
	}

	err = requests.ConnectionToSourcesService()
	for err != nil {
		log.Fatalf("Error connecting to sources service: %s\n", err.Error())
	}

	r := gin.Default()

	r.PUT("/order/assign", controllers.AssignOrderRequestHandler)
	r.GET("/order/acquire", controllers.AcquireOrderRequestHandler)
	r.POST("/order/cancel", controllers.CancelOrderRequestHandler)

	testing := r.Group("/testing")
	{
		testing.POST("/clean-database", controllers.CleanDatabaseRequestHandler)
		testing.POST("/clean-test-orders", controllers.CleanTestOrdersHandler)
	}

	listeningLine := ":8080"
	log.Printf("listening at %s\n", listeningLine)
	r.Run(listeningLine)
}
