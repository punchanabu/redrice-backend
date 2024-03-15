package routers

import (
	"github.com/gin-gonic/gin"
	docs "github.com/punchanabu/redrice-backend-go/docs"
	"github.com/punchanabu/redrice-backend-go/routers/api"
	"github.com/punchanabu/redrice-backend-go/routers/api/v1/"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func UseRouter() *gin.Engine {
	r := gin.New()
	docs.SwaggerInfo.Title = "RedRice API"
	docs.SwaggerInfo.Description = "This is a server for managing restaurant with RedRice API."
	docs.SwaggerInfo.BasePath = "/api/v1"
	apiv1 := r.Group("/api/v1")
	apiv1.Use()
	{

	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
