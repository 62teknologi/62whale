package main

import (
	"fmt"
	"net/http"
	"whale/62teknologi-golang-utility/utils"
	"whale/app/http/controllers"
	"whale/app/http/middlewares"
	"whale/app/interfaces"
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
		RegisterRoute(apiV1, "comment", &controllers.CommentController{})
		RegisterRoute(apiV1, "category", &controllers.CategoryController{})
		RegisterRoute(apiV1, "catalog", &controllers.CatalogController{})
		RegisterRoute(apiV1, "group", &controllers.GroupController{})
		RegisterRoute(apiV1, "item", &controllers.ItemController{})
		RegisterRoute(apiV1, "review", &controllers.ReviewController{})
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

func RegisterRoute(r gin.IRoutes, t string, c interfaces.Crud) {
	r.GET("/"+t+"/:table/:id", c.Find)
	r.GET("/"+t+"/:table/slug/:slug", c.Find)
	r.GET("/"+t+"/:table", c.FindAll)
	r.POST("/"+t+"/:table", c.Create)
	r.PUT("/"+t+"/:table/:id", c.Update)
	r.DELETE("/"+t+"/:table/:id", c.Delete)
}
