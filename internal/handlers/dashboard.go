package handlers

import (
	"casion/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	userID := c.GetString("user_id")
	var user models.User
	if result := h.db.First(&user, "user_id = ?", userID); result.Error != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  "error",
			Message: "User not found",
		})
		return
	}

	var totalTransactions int64
	if result := h.db.Model(&models.Transaction{}).Where("user_id = ?", userID).Count(&totalTransactions); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch transaction count",
		})
		return
	}

	var totalTransferred float64
	if result := h.db.Model(&models.Transaction{}).Where("user_id = ? AND type = ?", userID, "transfer").Select("COALESCE(SUM(amount), 0)").Scan(&totalTransferred); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to calculate total transferred amount",
		})
		return
	}

	var totalReceived float64
	if result := h.db.Model(&models.Transaction{}).Where("target_user_id = ? AND type = ?", userID, "transfer").Select("COALESCE(SUM(amount), 0)").Scan(&totalReceived); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to calculate total received amount",
		})
		return
	}

	c.JSON(http.StatusOK, DashboardStatsResponse{
		Status: "success",
		Result: DashboardStats{
			Balance:           user.Balance,
			TotalTransactions: totalTransactions,
			TotalTransferred:  totalTransferred,
			TotalReceived:     totalReceived,
		},
	})
}

func (h *DashboardHandler) GetRecentTransfers(c *gin.Context) {
	userID := c.GetString("user_id")
	var transfers []models.Transaction
	if result := h.db.Where("(user_id = ? OR target_user_id = ?) AND type = ?", userID, userID, "transfer").Order("created_at desc").Limit(10).Find(&transfers); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch recent transfers",
		})
		return
	}

	var response []RecentTransfer
	for _, t := range transfers {
		transferType := "sent"
		if t.TargetUserID == userID {
			transferType = "received"
		}

		response = append(response, RecentTransfer{
			TransferID: t.TransactionID,
			Amount:     t.Amount,
			Type:       transferType,
			Status:     t.Status,
			CreatedAt:  t.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, RecentTransfersResponse{
		Status: "success",
		Result: response,
	})
}

func (h *DashboardHandler) GetFailedTransfers(c *gin.Context) {
	userID := c.GetString("user_id")
	var transfers []models.Transaction
	if result := h.db.Where("user_id = ? AND type = ? AND status = ?", userID, "transfer", "failed").Order("created_at desc").Find(&transfers); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch failed transfers",
		})
		return
	}

	var response []FailedTransfer
	for _, t := range transfers {
		response = append(response, FailedTransfer{
			TransferID: t.TransactionID,
			Amount:     t.Amount,
			Remarks:    t.Remarks,
			CreatedAt:  t.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, FailedTransfersResponse{
		Status: "success",
		Result: response,
	})
}

type DashboardStatsResponse struct {
	Status string         `json:"status"`
	Result DashboardStats `json:"result"`
}

type DashboardStats struct {
	Balance           float64 `json:"balance"`
	TotalTransactions int64   `json:"total_transactions"`
	TotalTransferred  float64 `json:"total_transferred"`
	TotalReceived     float64 `json:"total_received"`
}

type RecentTransfersResponse struct {
	Status string           `json:"status"`
	Result []RecentTransfer `json:"result"`
}

type RecentTransfer struct {
	TransferID string    `json:"transfer_id"`
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_date"`
}

type FailedTransfersResponse struct {
	Status string           `json:"status"`
	Result []FailedTransfer `json:"result"`
}

type FailedTransfer struct {
	TransferID string    `json:"transfer_id"`
	Amount     float64   `json:"amount"`
	Remarks    string    `json:"remarks"`
	CreatedAt  time.Time `json:"created_date"`
}
