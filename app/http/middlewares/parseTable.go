package middlewares

import (
	"fmt"
	"whale/62teknologi-golang-utility/utils"

	"github.com/gin-gonic/gin"
)

// to do : set singular and plural table from json
func ParseTableMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("table:")
		fmt.Println(ctx.Param("table"))
		utils.SetPluralizeNames(ctx.Param("table"))
		ctx.Next()
	}
}
