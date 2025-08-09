package controller

import (
	"Student-Assistant-App/src/dtos/request"
	"Student-Assistant-App/src/dtos/response"
	"Student-Assistant-App/src/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService  service.UserService
	authService  service.AuthService
	otpService   service.OTPService
	emailService service.EmailService
}

func NewUserController(userService service.UserService, authService service.AuthService, otpService service.OTPService, emailService service.EmailService) *UserController {
	return &UserController{
		userService:  userService,
		authService:  authService,
		otpService:   otpService,
		emailService: emailService,
	}
}

// Send OTP endpoint
func (uc *UserController) SendOTP(ctx *gin.Context) {
	var sendOTPRequest request.SendOTPRequest
	err := ctx.ShouldBindJSON(&sendOTPRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	// Validate purpose
	if sendOTPRequest.Purpose != "signup" && sendOTPRequest.Purpose != "login" && sendOTPRequest.Purpose != "password_reset" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid OTP purpose"})
		return
	}

	// For signup, check if user doesn't exist
	if sendOTPRequest.Purpose == "signup" {
		existingUser, _ := uc.userService.GetUserByEmail(ctx.Request.Context(), sendOTPRequest.Email)
		if existingUser != nil {
			ctx.JSON(http.StatusConflict, gin.H{"message": "User already exists with this email"})
			return
		}
	}

	// For login, check if user exists
	if sendOTPRequest.Purpose == "login" || sendOTPRequest.Purpose == "password_reset" {
		existingUser, _ := uc.userService.GetUserByEmail(ctx.Request.Context(), sendOTPRequest.Email)
		if existingUser == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found with this email"})
			return
		}
	}

	err = uc.otpService.GenerateAndSendOTP(ctx.Request.Context(), sendOTPRequest.Email, sendOTPRequest.Purpose)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send OTP: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// Verify OTP endpoint
func (uc *UserController) VerifyOTP(ctx *gin.Context) {
	var verifyOTPRequest request.VerifyOTPRequest
	err := ctx.ShouldBindJSON(&verifyOTPRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	err = uc.otpService.VerifyOTP(ctx.Request.Context(), verifyOTPRequest.Email, verifyOTPRequest.Code, verifyOTPRequest.Purpose)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}

// Resend OTP endpoint
func (uc *UserController) ResendOTP(ctx *gin.Context) {
	var sendOTPRequest request.SendOTPRequest
	err := ctx.ShouldBindJSON(&sendOTPRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	err = uc.otpService.ResendOTP(ctx.Request.Context(), sendOTPRequest.Email, sendOTPRequest.Purpose)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "OTP resent successfully"})
}

// Signup with OTP verification
func (uc *UserController) SignupWithOTP(ctx *gin.Context) {
	var signupRequest request.SignupWithOTPRequest
	err := ctx.ShouldBindJSON(&signupRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.CreateUserResponse{
			Message: "Bad request",
		})
		return
	}

	if signupRequest.GetName() == "" || signupRequest.GetEmail() == "" || signupRequest.GetPassword() == "" {
		ctx.JSON(http.StatusBadRequest, response.CreateUserResponse{
			Message: "Name, email and password are required",
		})
		return
	}

	// Verify OTP first
	err = uc.otpService.VerifyOTP(ctx.Request.Context(), signupRequest.Email, signupRequest.OTPCode, "signup")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.CreateUserResponse{
			Message: "Invalid or expired OTP",
		})
		return
	}

	// Create user
	createUserRequest := &request.CreateUserRequest{
		Name:     signupRequest.Name,
		Email:    signupRequest.Email,
		Password: signupRequest.Password,
		Role:     signupRequest.Role,
	}

	createUserResponse, err := uc.userService.CreateUser(ctx.Request.Context(), createUserRequest)
	if err != nil {
		statusCode := http.StatusBadRequest
		if strings.Contains(err.Error(), "already exists") {
			statusCode = http.StatusConflict
		}

		ctx.JSON(statusCode, response.CreateUserResponse{
			Message: err.Error(),
		})
		return
	}

	// Send welcome email
	uc.emailService.SendWelcomeEmail(signupRequest.Email, signupRequest.Name)

	ctx.JSON(http.StatusCreated, createUserResponse)
}

// Login with OTP
func (uc *UserController) LoginWithOTP(ctx *gin.Context) {
	var loginRequest request.LoginWithOTPRequest
	err := ctx.ShouldBindJSON(&loginRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	// Verify OTP
	err = uc.otpService.VerifyOTP(ctx.Request.Context(), loginRequest.Email, loginRequest.OTPCode, "login")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired OTP"})
		return
	}

	// Get user and generate token
	user, err := uc.userService.GetUserByEmail(ctx.Request.Context(), loginRequest.Email)
	if err != nil || user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	// Generate JWT token
	token, err := uc.authService.GenerateTokenForUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, response.LoginResponse{
		Message: "Login successful",
		User:    user,
		Token:   token,
	})
}

// Traditional Signup (without OTP - for backward compatibility)
func (uc *UserController) Signup(ctx *gin.Context) {
	var createUserRequest request.CreateUserRequest
	err := ctx.ShouldBindJSON(&createUserRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.CreateUserResponse{
			Message: "Bad request",
		})
		return
	}

	if createUserRequest.GetName() == "" || createUserRequest.GetEmail() == "" || createUserRequest.GetPassword() == "" {
		ctx.JSON(http.StatusBadRequest, response.CreateUserResponse{
			Message: "Name, email and password are required",
		})
		return
	}

	createUserResponse, err := uc.userService.CreateUser(ctx.Request.Context(), &createUserRequest)
	if err != nil {
		statusCode := http.StatusBadRequest
		if strings.Contains(err.Error(), "already exists") {
			statusCode = http.StatusConflict
		}

		ctx.JSON(statusCode, response.CreateUserResponse{
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, createUserResponse)
}

// Traditional Login (without OTP - for backward compatibility)
func (uc *UserController) Login(ctx *gin.Context) {
	var loginRequest request.LoginRequest
	err := ctx.ShouldBindJSON(&loginRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.LoginResponse{
			Message: "Bad request",
		})
		return
	}

	loginResponse, err := uc.authService.Login(ctx.Request.Context(), &loginRequest)
	if err != nil {
		statusCode := http.StatusUnauthorized
		ctx.JSON(statusCode, response.LoginResponse{
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, loginResponse)
}

// Get user by ID
func (uc *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "User ID is required"})
		return
	}

	user, err := uc.userService.GetUserByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User retrieved successfully",
		"user":    user,
	})
}

// Get all users
func (uc *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := uc.userService.GetAllUsers(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve users"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"users":   users,
	})
}

// Update user
func (uc *UserController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "User ID is required"})
		return
	}

	var updateUserRequest request.UpdateUserRequest
	err := ctx.ShouldBindJSON(&updateUserRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	user, err := uc.userService.UpdateUser(ctx.Request.Context(), id, &updateUserRequest)
	if err != nil {
		statusCode := http.StatusNotFound
		if strings.Contains(err.Error(), "validation") {
			statusCode = http.StatusBadRequest
		}
		ctx.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}

// Delete user
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "User ID is required"})
		return
	}

	err := uc.userService.DeleteUser(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Get current user (from JWT token)
func (uc *UserController) GetCurrentUser(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}

	user, err := uc.userService.GetUserByID(ctx.Request.Context(), userID.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Current user retrieved successfully",
		"user":    user,
	})
}