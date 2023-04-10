package main

import (
	"fmt"
	"whale/app/http/controllers"
	middlewares "whale/app/http/midlewares"
	"whale/config"
	"whale/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("cannot load config: %w", err)
		return
	}

	utils.ConnectDatabase(config)

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
