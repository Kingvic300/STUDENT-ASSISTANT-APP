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

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not set in .env file")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME not set in .env file")
	}
	db := client.Database(dbName)

	userRepo := repository.NewUserRepositoryImpl(db)
	otpRepo := repository.NewOTPRepositoryImpl(db)

	emailService, err := service.NewEmailService()
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}

	otpService := service.NewOTPService(otpRepo, emailService)
	userService := service.NewUserServiceImpl(userRepo)
	authService := service.NewAuthService(userService)

	userController := controller.NewUserController(userService, authService, otpService, emailService)

	router := gin.Default()

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

	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/users/me", userController.GetCurrentUser)
		api.GET("/users/:id", userController.GetUser)
		api.PUT("/users/:id", userController.UpdateUser)
		api.DELETE("/users/:id", userController.DeleteUser)

		admin := api.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/users", userController.GetAllUsers)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		log.Printf("Starting server on :%s", port)
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received, exiting...")
}
