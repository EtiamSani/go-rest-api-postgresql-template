package route

import (
	"github.com/etiamsani/go-rest-api-postgresl-template/api/handler"
	"github.com/gin-gonic/gin"
)

func UseRouter() *gin.Engine {
router := gin.Default()
router.POST("/use", handler.CreateUser) 
return router
}