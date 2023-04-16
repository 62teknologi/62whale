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

	utils.InitPluralize()

	r := gin.Default()

	apiV1 := r.Group("/whale/v1").Use(middlewares.DbSelectorMiddleware())
	{
		apiV1.GET("/products", controllers.FindProducts)
		apiV1.POST("/products", controllers.CreateProduct)
		apiV1.GET("/products/:id", controllers.FindProduct)
		apiV1.PUT("/products/:id", controllers.UpdateProduct)
		//apiV1.DELETE("/users/:id", controllers.DeleteProduct)

		apiV1.GET("/catalog/:table", controllers.FindCatalogues).Use(middlewares.ParseTableMiddleware())
		apiV1.POST("/catalog/:table", controllers.CreateCatalog).Use(middlewares.ParseTableMiddleware())
		apiV1.GET("/catalog/:table/:id", controllers.FindCatalog).Use(middlewares.ParseTableMiddleware())
		apiV1.PUT("/catalog/:table/:id", controllers.UpdateCatalog).Use(middlewares.ParseTableMiddleware())
		apiV1.DELETE("/catalog/:table/:id", controllers.DeleteCatalog).Use(middlewares.ParseTableMiddleware())
	}

	err = r.Run(config.HTTPServerAddress)

	if err != nil {
		fmt.Printf("cannot run server: %w", err)
		return
	}
}
