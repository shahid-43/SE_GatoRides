package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/verify-email", controllers.VerifyEmail)

	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/user/provide-ride", controllers.ProvideRide)
		protected.POST("/user/location", controllers.UpdateUserLocation)
		protected.POST("/user/logout", controllers.LogOut)
		protected.GET("/home", controllers.HomeHandler)

		// protected.GET("/rides/feed", controllers.FetchRideFeed)

	}

	return r
}
