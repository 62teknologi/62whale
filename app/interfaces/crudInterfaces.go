package interfaces

import "github.com/gin-gonic/gin"

type Crud interface {
	Find(*gin.Context)
	FindAll(*gin.Context)
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	DeleteByQuery(*gin.Context)
}
