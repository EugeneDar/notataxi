package main

import (
	"app/src/services/orders/app/controllers"
	"app/src/services/orders/app/database"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	err := database.EstablishConnection()
	for err != nil {
		log.Println(err, "Error connecting to database: %s; retrying...", err.Error())
		err = database.EstablishConnection()
	}

	r := gin.Default()

	service, err := controllers.NewService()
	if err != nil {
		log.Fatal(err)
	}

	r.PUT("/order/assign", service.AssignOrderRequestHandler)
	r.GET("/order/acquire", service.AcquireOrderRequestHandler)
	r.POST("/order/cancel", service.CancelOrderRequestHandler)

	testing := r.Group("/testing")
	{
		testing.POST("/clean-database", service.CleanDatabaseRequestHandler)
		testing.POST("/clean-test-orders", service.CleanTestOrdersHandler)
	}

	listeningLine := ":8080"
	log.Printf("listening at %s\n", listeningLine)
	r.Run(listeningLine)
}
