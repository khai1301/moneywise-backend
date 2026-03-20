package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/khai1301/moneywise-backend/internal/config"
	"github.com/khai1301/moneywise-backend/internal/handler"
	"github.com/khai1301/moneywise-backend/internal/middleware"
	"github.com/khai1301/moneywise-backend/internal/repository"
	"github.com/khai1301/moneywise-backend/internal/service"
)

func SetupRoutes(router *gin.Engine) {
	// CORS Configuration (Quan trọng: Cho phép Frontend truy cập)
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://moneywise-nu.vercel.app",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Setup Authentication Module Dependencies
	userRepo := repository.NewUserRepository(config.DB)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// Setup Category Module Dependencies
	categoryRepo := repository.NewCategoryRepository(config.DB)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Setup Transaction Module Dependencies
	transactionRepo := repository.NewTransactionRepository(config.DB)
	transactionService := service.NewTransactionService(transactionRepo, categoryRepo)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// Setup Budget Module Dependencies
	budgetRepo := repository.NewBudgetRepository(config.DB)
	budgetService := service.NewBudgetService(budgetRepo, categoryRepo, config.DB)
	budgetHandler := handler.NewBudgetHandler(budgetService)

	// Setup Analytics Module Dependencies
	analyticsHandler := handler.NewAnalyticsHandler()

	api := router.Group("/api")
	{
		// Auth routes
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		// Protected Routes
		protectedGroup := api.Group("")
		protectedGroup.Use(middleware.RequireAuth())
		{
			// Category routes
			categoryGroup := protectedGroup.Group("/categories")
			{
				categoryGroup.POST("", categoryHandler.Create)
				categoryGroup.GET("", categoryHandler.GetAll)
				categoryGroup.PUT("/:id", categoryHandler.Update)
				categoryGroup.DELETE("/:id", categoryHandler.Delete)
			}

			// Budget routes
			budgetGroup := protectedGroup.Group("/budgets")
			{
				budgetGroup.POST("", budgetHandler.Create)
				budgetGroup.GET("", budgetHandler.GetAll)
				budgetGroup.PUT("/:id", budgetHandler.Update)
				budgetGroup.DELETE("/:id", budgetHandler.Delete)
			}

			// Analytics routes
			analyticsGroup := protectedGroup.Group("/analytics")
			{
				analyticsGroup.GET("/monthly", analyticsHandler.Monthly)
				analyticsGroup.GET("/categories", analyticsHandler.CategorySummary)
			}

			// Transaction routes
			transactionGroup := protectedGroup.Group("/transactions")
			{
				transactionGroup.POST("", transactionHandler.Create)
				transactionGroup.GET("", transactionHandler.GetAll)
				transactionGroup.GET("/:id", transactionHandler.GetByID)
				transactionGroup.PUT("/:id", transactionHandler.Update)
				transactionGroup.DELETE("/:id", transactionHandler.Delete)
			}
		}
	}
}
