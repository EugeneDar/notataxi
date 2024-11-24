package main

import (
	"app/src/services/orders/app/controllers"
	"app/src/services/orders/app/database"
	"app/src/services/orders/app/utils"
	"fmt"
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
	r.PUT("/order/assign", controllers.AssignOrderRequestHandler)
	r.GET("/order/acquire", controllers.AcquireOrderRequestHandler)
	r.POST("/order/cancel", controllers.CancelOrderRequestHandler)

	testing := r.Group("/testing")
    {
		testing.POST("/clean-database", controllers.CleanDatabaseRequestHandler)
        testing.POST("/clean-test-orders", controllers.CleanTestOrdersHandler)
    }

	listeningLine := fmt.Sprintf(":%s", utils.GetenvSafe("ORDERS_SERVICE_PORT"))
	log.Printf("listening at %s\n", listeningLine)
	r.Run(listeningLine)
}
