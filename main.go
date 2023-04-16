package main

import (
	"fmt"
	"net/http"
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

	apiV1 := r.Group("/api/v1").Use(middlewares.DbSelectorMiddleware())
	{
		apiV1.GET("/category/:table", controllers.FindCatalogCategories).Use(middlewares.ParseTableMiddleware())
		apiV1.POST("/category/:table", controllers.CreateCatalogCategory).Use(middlewares.ParseTableMiddleware())
		apiV1.GET("/category/:table/:id", controllers.FindCatalogCategory).Use(middlewares.ParseTableMiddleware())
		apiV1.PUT("/category/:table/:id", controllers.UpdateCatalogCategory).Use(middlewares.ParseTableMiddleware())
		apiV1.DELETE("/category/:table/:id", controllers.DeleteCatalogCategory).Use(middlewares.ParseTableMiddleware())

		apiV1.GET("/group/:table", controllers.FindCatalogGroups).Use(middlewares.ParseTableMiddleware())
		apiV1.POST("/group/:table", controllers.CreateCatalogGroup).Use(middlewares.ParseTableMiddleware())
		apiV1.GET("/group/:table/:id", controllers.FindCatalogGroup).Use(middlewares.ParseTableMiddleware())
		apiV1.PUT("/group/:table/:id", controllers.UpdateCatalogGroup).Use(middlewares.ParseTableMiddleware())
		apiV1.DELETE("/group/:table/:id", controllers.DeleteCatalogGroup).Use(middlewares.ParseTableMiddleware())

		apiV1.GET("/catalog/:table", controllers.FindCatalogues).Use(middlewares.ParseTableMiddleware())
		apiV1.POST("/catalog/:table", controllers.CreateCatalog).Use(middlewares.ParseTableMiddleware())
		apiV1.GET("/catalog/:table/:id", controllers.FindCatalog).Use(middlewares.ParseTableMiddleware())
		apiV1.PUT("/catalog/:table/:id", controllers.UpdateCatalog).Use(middlewares.ParseTableMiddleware())
		apiV1.DELETE("/catalog/:table/:id", controllers.DeleteCatalog).Use(middlewares.ParseTableMiddleware())
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, utils.ResponseData("success", "Server running well", nil))
	})

	err = r.Run(config.HTTPServerAddress)

	if err != nil {
		fmt.Printf("cannot run server: %w", err)
		return
	}
}
