package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hendraadwii/library/cmd/api/handlers"
	"github.com/hendraadwii/library/internal/auth"
	"github.com/hendraadwii/library/internal/config"
	"github.com/hendraadwii/library/internal/database"
	"github.com/hendraadwii/library/internal/models"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Convert config types
func toDBConfig(cfg config.DBConfig) database.DBConfig {
	return database.DBConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		DBName:   cfg.DBName,
	}
}

func toAuthConfig(cfg config.AuthConfig) auth.Config {
	return auth.Config{
		AccessTokenSecret:  cfg.AccessTokenSecret,
		RefreshTokenSecret: cfg.RefreshTokenSecret,
		AccessTokenExpiry:  cfg.AccessTokenExpiry,
		RefreshTokenExpiry: cfg.RefreshTokenExpiry,
		Issuer:             cfg.Issuer,
	}
}

// Run starts the server
func Run() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create a new Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Connect to database
	db, err := database.NewDBConnection(toDBConfig(cfg.Database))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize services
	bookService := models.NewBookService(db)
	userService := models.NewUserService(db)
	borrowingService := models.NewBorrowingService(db)
	passwordResetService := models.NewPasswordResetService(db)

	// Initialize auth components
	tokenManager := auth.NewTokenManager(toAuthConfig(cfg.Auth))
	authMiddleware := auth.NewMiddleware(tokenManager)
	rateLimiter := auth.NewRateLimiter(
		cfg.RateLimiting.MaxRequests,
		cfg.RateLimiting.PerDuration,
		cfg.RateLimiting.CleanupEvery,
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, tokenManager, rateLimiter)
	authHandler.PasswordResetService = passwordResetService
	bookHandler := handlers.NewBookHandler(bookService)
	borrowingHandler := handlers.NewBorrowingHandler(borrowingService, bookService)
	userHandler := handlers.NewUserHandler(userService)

	// Group routes with API version
	api := router.Group("/api/v1")

	// Auth routes - with rate limiting
	authRoutes := api.Group("/auth")
	authRoutes.Use(rateLimiter.RateLimit())
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/refresh", authHandler.RefreshToken)
		authRoutes.POST("/forgot-password", authHandler.ForgotPassword)
		authRoutes.POST("/verify-pin", authHandler.VerifyPin)
		authRoutes.POST("/reset-password", authHandler.ResetPassword)
	}

	// Book routes - some with authentication
	bookRoutes := api.Group("/books")
	{
		// Public routes
		bookRoutes.GET("", bookHandler.GetBooks)
		bookRoutes.GET("/:id", bookHandler.GetBook)
		bookRoutes.GET("/most-borrowed", borrowingHandler.GetMostBorrowedBooks)

		// Protected routes (require authentication)
		bookRoutes.Use(authMiddleware.AuthRequired())
		bookRoutes.POST("", bookHandler.CreateBook)
		bookRoutes.PUT("/:id", bookHandler.UpdateBook)
		bookRoutes.DELETE("/:id", bookHandler.DeleteBook)
	}

	// Borrowing routes - all protected
	borrowingRoutes := api.Group("/borrowings")
	borrowingRoutes.Use(authMiddleware.AuthRequired())
	{
		borrowingRoutes.GET("", borrowingHandler.GetBorrowings)
		borrowingRoutes.GET("/:id", borrowingHandler.GetBorrowing)
		borrowingRoutes.POST("", borrowingHandler.BorrowBook)
		borrowingRoutes.POST("/:id/return", borrowingHandler.ReturnBook)
		borrowingRoutes.GET("/overdue", borrowingHandler.GetOverdueBooks)
		borrowingRoutes.GET("/member", borrowingHandler.GetMemberBorrowings)
	}

	// User routes - admin only
	userRoutes := api.Group("/users")
	userRoutes.Use(authMiddleware.AuthRequired(), authMiddleware.RoleRequired(string(models.RoleAdmin)))
	{
		userRoutes.GET("", userHandler.GetUsers)
		userRoutes.POST("", userHandler.CreateUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
		userRoutes.PUT("/:id/password", userHandler.ResetUserPassword)
		userRoutes.PUT("/:id/status", userHandler.ToggleUserStatus)
		userRoutes.GET("/count", userHandler.GetUserCount)
	}

	// Serve static files for book covers
	router.Static("/static/cover", "./uploads/cover")

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:         ":8081",
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Run the server in a goroutine
	go func() {
		log.Printf("Server starting on port 8081...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped gracefully")
}
