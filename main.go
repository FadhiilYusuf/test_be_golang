package main

import (
	"log"

	"github.com/cngJo/golang-api-auth/controllers"
	"github.com/cngJo/golang-api-auth/database"
	"github.com/cngJo/golang-api-auth/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()
	database.Migrate()

	router := initRouter()
	router.Run(":8080")
}

func initRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.RegisterUser)
			auth.POST("/login", controllers.LoginUser)
			auth.POST("/reset-password", controllers.ResetPassword)
			auth.POST("/logout", controllers.LogoutUser)
			auth.POST("/profile", controllers.UpdateProfile)
			auth.POST("/products", controllers.CreateProduct)
			auth.PUT("/products/:id", controllers.UpdateProduct)
			auth.DELETE("/products/:id", controllers.DeleteProduct)
			auth.GET("/products/:id", controllers.GetProductDetail)
			auth.GET("/orders", controllers.GetOrders)
			auth.POST("/orders", controllers.CreateOrderStatus)

		}

		secured := api.Group("").Use(middlewares.Auth())
		{
			secured.GET("/ping", controllers.Ping)
		}
	}

	return router
}
