package main

import (
	"casion/internal/config"
	"casion/internal/handlers"
	"casion/internal/middleware"
	"casion/internal/worker"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db := config.InitDB()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	profileHandler := handlers.NewProfileHandler(db)
	transactionHandler := handlers.NewTransactionHandler(db)
	dashboardHandler := handlers.NewDashboardHandler(db)

	// Initialize router
	router := gin.Default()

	// Serve OpenAPI specification
	router.StaticFile("/openapi.yaml", "./api/openapi.yaml")
	router.Static("/swagger-ui", "./api/swagger-ui")

	// API routes
	api := router.Group("/api")
	{
		// Auth routes
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Profile routes
			protected.PUT("/profile", profileHandler.UpdateProfile)

			// Transaction routes
			protected.POST("/topup", transactionHandler.TopUp)
			protected.POST("/payment", transactionHandler.Payment)
			protected.POST("/transfer", transactionHandler.Transfer)
			protected.GET("/transactions", transactionHandler.GetTransactions)

			// Dashboard routes
			protected.GET("/dashboard/stats", dashboardHandler.GetDashboardStats)
			protected.GET("/dashboard/transfers/recent", dashboardHandler.GetRecentTransfers)
			protected.GET("/dashboard/transfers/failed", dashboardHandler.GetFailedTransfers)
		}
	}

	// Start transfer worker
	go worker.ProcessTransfers(db)

	// Start server
	log.Fatal(router.Run(":8080"))
}
