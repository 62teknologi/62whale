package middlewares

import (
	"github.com/62teknologi/62whale/62golib/utils"

	"github.com/gin-gonic/gin"
)

func DbSelectorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if ctx.Query("db") == "2" {
			utils.DB = utils.DB2
		} else {
			utils.DB = utils.DB1
		}
		ctx.Next()
	}
}
