package middlewares

import (
	"whale/62teknologi-golang-utility/utils"

	"github.com/gin-gonic/gin"
)

// to do : set singular and plural table from json
func ParseTableMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		utils.SetPluralizeNames(ctx.Param("table"))
		ctx.Next()
	}
}
