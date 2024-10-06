package main

import (
	"fabric-rest-api/controller"
	"fmt"

	"github.com/gin-gonic/gin"
)


func main() {
	r := gin.Default()

	// Define routes
	r.POST("/createProduct", controller.CreateProductHandler)
	r.POST("/supplyProduct", controller.SupplyProductHandler)
	r.POST("/wholesaleProduct", controller.WholesaleProductHandler)
	r.GET("/queryProduct", controller.QueryProductHandler)
	r.POST("/sellProduct", controller.SellProductHandler)

	// Start the server
	fmt.Println("Server is running on http://localhost:3000")
	r.Run("localhost:3000")
}
