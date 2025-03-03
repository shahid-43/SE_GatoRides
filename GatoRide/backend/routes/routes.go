package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.POST("/signup", gin.WrapF(controllers.Signup))
	r.POST("/login", gin.WrapF(controllers.Login))
	r.GET("/verify-email", gin.WrapF(controllers.VerifyEmail))

	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/home", controllers.HomeHandler)
		protected.GET("/rides/feed", controllers.FetchRideFeed)
		protected.POST("/user/provide-ride", controllers.ProvideRide)
		protected.POST("/users/location", controllers.UpdateUserLocation)
	}

	return r
}
