package main

import (
	"fmt"
	"whale/62teknologi-golang-utility/utils"
	"whale/app/http/controllers"
	"whale/app/http/middlewares"
	"whale/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("cannot load config: %w", err)
		return
	}

	// todo : replace last variable with spread notation "..."
	utils.ConnectDatabase(config.DBDriver, config.DBSource1, config.DBSource2)

	r := gin.Default()

	apiV1 := r.Group("/api/v1").Use(middlewares.DbSelectorMiddleware())
	{
		apiV1.GET("/products", controllers.FindProducts)
		apiV1.POST("/products", controllers.CreateProduct)
		apiV1.GET("/products/:id", controllers.FindProduct)
		apiV1.PUT("/products/:id", controllers.UpdateProduct)
		//apiV1.DELETE("/users/:id", controllers.DeleteProduct)
	}

	err = r.Run(config.HTTPServerAddress)

	if err != nil {
		fmt.Printf("cannot run server: %w", err)
		return
	}
}
