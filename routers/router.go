package routers

import (
	"github.com/gin-gonic/gin"
	docs "github.com/punchanabu/redrice-backend-go/docs"
	"github.com/punchanabu/redrice-backend-go/middleware"
	"github.com/punchanabu/redrice-backend-go/routers/api"
	v1 "github.com/punchanabu/redrice-backend-go/routers/api/v1"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func UseRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())

	docs.SwaggerInfo.Title = "RedRice API"
	docs.SwaggerInfo.Description = "This is a server for managing restaurant with RedRice API build with Go Gin and Gorm"
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	apiv1 := r.Group("/api/v1")
	auth := apiv1.Group("/auth")
	auth.POST("/login", api.Login)
	auth.POST("/register", api.Register)
	apiv1.Use(middleware.Auth())
	{
		// for authorized user
		apiv1.GET("/restaurants", v1.GetRestaurants)
		apiv1.GET("/restaurants/:id", v1.GetRestaurant)
		apiv1.GET("/reservations", v1.GetReservations)
		apiv1.GET("/reservations/:id", v1.GetReservation)
		apiv1.GET("/users", v1.GetUsers)
		apiv1.GET("/users/:id", v1.GetUser)
		apiv1.POST("/reservations", v1.CreateReservation)

		// for admin
		adminRoutes := apiv1.Group("/")
		adminRoutes.Use(middleware.Admin())
		{
			adminRoutes.POST("/users", v1.CreateUser)
			adminRoutes.PUT("/users/:id", v1.UpdateUser)
			adminRoutes.DELETE("/users/:id", v1.DeleteUser)
			adminRoutes.POST("/restaurants", v1.CreateRestaurant)
			adminRoutes.PUT("/reservations/:id", v1.UpdateReservation)
			adminRoutes.PUT("/restaurants/:id", v1.UpdateRestaurant)
			adminRoutes.DELETE("/restaurants/:id", v1.DeleteRestaurant)
		}

	}
	return r
}
