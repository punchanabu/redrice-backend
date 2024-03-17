package routers

import (
	"github.com/gin-gonic/gin"
	docs "github.com/punchanabu/redrice-backend-go/docs"
	"github.com/punchanabu/redrice-backend-go/routers/api"
	v1 "github.com/punchanabu/redrice-backend-go/routers/api/v1"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func UseRouter() *gin.Engine {
	r := gin.New()

	docs.SwaggerInfo.Title = "RedRice API"
	docs.SwaggerInfo.Description = "This is a server for managing restaurant with RedRice API."
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/auth", api.Authenticate)
	apiv1 := r.Group("/api/v1")
	apiv1.Use()
	{
		// users
		apiv1.GET("/users", v1.GetUsers)
		apiv1.POST("/users", v1.CreateUser)
		apiv1.GET("/users/:id", v1.GetUser)
		apiv1.PUT("/users/:id", v1.UpdateUser)
		apiv1.DELETE("/users/:id", v1.DeleteUser)

		// restaurants
		apiv1.GET("/restaurants", v1.GetRestaurants)
		apiv1.POST("/restaurants", v1.CreateRestaurant)
		apiv1.GET("/restaurants/:id", v1.GetRestaurant)
		apiv1.PUT("/restaurants/:id", v1.UpdateRestaurant)
		apiv1.DELETE("/restaurants/:id", v1.DeleteRestaurant)

		// reservation
		apiv1.GET("/reservations", v1.GetReservations)
		apiv1.POST("/reservations", v1.CreateReservation)
		apiv1.GET("/reservations/:id", v1.GetReservation)
		apiv1.PUT("/reservations/:id", v1.UpdateReservation)
	}
	return r
}
