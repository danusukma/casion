package handlers

import (
	"casion/internal/models"
	"casion/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid request payload",
		})
		return
	}

	// Check if phone number already exists
	var existingUser models.User
	if result := h.db.Where("phone_number = ?", request.PhoneNumber).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Status:  "error",
			Message: "Phone number already exists",
		})
		return
	}

	// Create new user
	user := models.User{
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		PhoneNumber: request.PhoneNumber,
		Address:     request.Address,
		Pin:         request.Pin,
	}

	if result := h.db.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{
		Status: "success",
		Result: UserResponse{
			UserID:      user.UserID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			PhoneNumber: user.PhoneNumber,
			Address:     user.Address,
			Balance:     user.Balance,
			CreatedAt:   user.CreatedAt,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid request payload",
		})
		return
	}

	// Find user by phone number
	var user models.User
	if result := h.db.Where("phone_number = ?", request.PhoneNumber).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid phone number or PIN",
		})
		return
	}

	// Verify PIN
	if user.Pin != request.Pin {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid phone number or PIN",
		})
		return
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to generate access token",
		})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to generate refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Status: "success",
		Result: TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

type RegisterRequest struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Pin         string `json:"pin" binding:"required"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Pin         string `json:"pin" binding:"required"`
}

type RegisterResponse struct {
	Status string       `json:"status"`
	Result UserResponse `json:"result"`
}

type LoginResponse struct {
	Status string        `json:"status"`
	Result TokenResponse `json:"result"`
}

type UserResponse struct {
	UserID      string    `json:"user_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	Balance     float64   `json:"balance"`
	CreatedAt   time.Time `json:"created_date"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
