package main

import (
	"Student-Assistant-App/src/controller"
	"Student-Assistant-App/src/data/repository"
	"Student-Assistant-App/src/middleware"
	"Student-Assistant-App/src/service"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// MongoDB connection URI
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not set in .env file")
	}

	// Create MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Context with timeout for initial connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ensure disconnect on shutdown
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()

	// Get DB name from env
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME not set in .env file")
	}
	db := client.Database(dbName)

	// Initialize repositories
	userRepo := repository.NewUserRepositoryImpl(db)
	otpRepo := repository.NewOTPRepositoryImpl(db)

	// Initialize email service with error handling
	emailService, err := service.NewEmailService()
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}

	// Initialize other services
	otpService := service.NewOTPService(otpRepo, emailService)
	userService := service.NewUserServiceImpl(userRepo)
	authService := service.NewAuthService(userService)

	// Initialize controllers
	userController := controller.NewUserController(userService, authService, otpService, emailService)

	// Setup Gin router
	router := gin.Default()

	// Public routes
	public := router.Group("/api")
	{
		public.POST("/auth/signup", userController.Signup)
		public.POST("/auth/login", userController.Login)

		public.POST("/auth/send-otp", userController.SendOTP)
		public.POST("/auth/verify-otp", userController.VerifyOTP)
		public.POST("/auth/resend-otp", userController.ResendOTP)
		public.POST("/auth/signup-with-otp", userController.SignupWithOTP)
		public.POST("/auth/login-with-otp", userController.LoginWithOTP)
	}

	// Protected routes with auth middleware
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/users/me", userController.GetCurrentUser)
		api.GET("/users/:id", userController.GetUser)
		api.PUT("/users/:id", userController.UpdateUser)
		api.DELETE("/users/:id", userController.DeleteUser)

		// Admin only routes
		admin := api.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/users", userController.GetAllUsers)
		}
	}

	// Server port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Run server in goroutine so we can listen for shutdown signals
	go func() {
		log.Printf("Starting server on :%s", port)
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received, exiting...")
}
