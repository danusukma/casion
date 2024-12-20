package handlers

import (
	"casion/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProfileHandler struct {
	db *gorm.DB
}

func NewProfileHandler(db *gorm.DB) *ProfileHandler {
	return &ProfileHandler{db: db}
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var request UpdateProfileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid request payload",
		})
		return
	}

	userID := c.GetString("user_id")
	var user models.User
	if result := h.db.First(&user, "user_id = ?", userID); result.Error != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  "error",
			Message: "User not found",
		})
		return
	}

	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.Address = request.Address

	if result := h.db.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to update profile",
		})
		return
	}

	c.JSON(http.StatusOK, UpdateProfileResponse{
		Status: "success",
		Result: ProfileResult{
			UserID:    user.UserID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Address:   user.Address,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Address   string `json:"address" binding:"required"`
}

type UpdateProfileResponse struct {
	Status string        `json:"status"`
	Result ProfileResult `json:"result"`
}

type ProfileResult struct {
	UserID    string    `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Address   string    `json:"address"`
	UpdatedAt time.Time `json:"updated_date"`
}
