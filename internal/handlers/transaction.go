package handlers

import (
	"casion/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	db *gorm.DB
}

func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{db: db}
}

func (h *TransactionHandler) TopUp(c *gin.Context) {
	var request TopUpRequest
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

	balanceBefore := user.Balance
	user.Balance += request.Amount

	tx := h.db.Begin()
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to update balance",
		})
		return
	}

	transaction := models.Transaction{
		TransactionID: uuid.New().String(),
		UserID:        userID,
		Type:          "top_up",
		Amount:        request.Amount,
		Status:        "success",
		BalanceBefore: balanceBefore,
		BalanceAfter:  user.Balance,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to create transaction record",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, TopUpResponse{
		Status: "success",
		Result: TopUpResult{
			TopUpID:       transaction.TransactionID,
			AmountTopUp:   transaction.Amount,
			BalanceBefore: transaction.BalanceBefore,
			BalanceAfter:  transaction.BalanceAfter,
			CreatedAt:     transaction.CreatedAt,
		},
	})
}

func (h *TransactionHandler) Payment(c *gin.Context) {
	var request PaymentRequest
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

	if user.Balance < request.Amount {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Insufficient balance",
		})
		return
	}

	balanceBefore := user.Balance
	user.Balance -= request.Amount

	tx := h.db.Begin()
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to update balance",
		})
		return
	}

	transaction := models.Transaction{
		TransactionID: uuid.New().String(),
		UserID:        userID,
		Type:          "payment",
		Amount:        request.Amount,
		Status:        "success",
		Remarks:       request.Remarks,
		BalanceBefore: balanceBefore,
		BalanceAfter:  user.Balance,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to create transaction record",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, PaymentResponse{
		Status: "success",
		Result: PaymentResult{
			PaymentID:     transaction.TransactionID,
			Amount:        transaction.Amount,
			Remarks:       transaction.Remarks,
			BalanceBefore: transaction.BalanceBefore,
			BalanceAfter:  transaction.BalanceAfter,
			CreatedAt:     transaction.CreatedAt,
		},
	})
}

func (h *TransactionHandler) Transfer(c *gin.Context) {
	var request TransferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid request payload",
		})
		return
	}

	userID := c.GetString("user_id")
	var sender models.User
	if result := h.db.First(&sender, "user_id = ?", userID); result.Error != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  "error",
			Message: "User not found",
		})
		return
	}

	if sender.Balance < request.Amount {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Insufficient balance",
		})
		return
	}

	var receiver models.User
	if result := h.db.First(&receiver, "phone_number = ?", request.TargetUser); result.Error != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Target user not found",
		})
		return
	}

	balanceBefore := sender.Balance
	sender.Balance -= request.Amount

	tx := h.db.Begin()
	if err := tx.Save(&sender).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to update sender balance",
		})
		return
	}

	transaction := models.Transaction{
		TransactionID: uuid.New().String(),
		UserID:        userID,
		TargetUserID:  receiver.UserID,
		Type:          "transfer",
		Amount:        request.Amount,
		Status:        "pending",
		Remarks:       request.Remarks,
		BalanceBefore: balanceBefore,
		BalanceAfter:  sender.Balance,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to create transaction record",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, TransferResponse{
		Status: "success",
		Result: TransferResult{
			TransferID:    transaction.TransactionID,
			Amount:        transaction.Amount,
			Remarks:       transaction.Remarks,
			BalanceBefore: transaction.BalanceBefore,
			BalanceAfter:  transaction.BalanceAfter,
			CreatedAt:     transaction.CreatedAt,
		},
	})
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID := c.GetString("user_id")
	var transactions []models.Transaction
	if result := h.db.Where("user_id = ?", userID).Order("created_at desc").Find(&transactions); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch transactions",
		})
		return
	}

	var response []TransactionResponse
	for _, t := range transactions {
		response = append(response, TransactionResponse{
			TransferID:      t.TransactionID,
			Status:          t.Status,
			UserID:          t.UserID,
			TransactionType: t.Type,
			Amount:          t.Amount,
			Remarks:         t.Remarks,
			BalanceBefore:   t.BalanceBefore,
			BalanceAfter:    t.BalanceAfter,
			CreatedAt:       t.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, GetTransactionsResponse{
		Status: "success",
		Result: response,
	})
}

type TopUpRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

type PaymentRequest struct {
	Amount  float64 `json:"amount" binding:"required"`
	Remarks string  `json:"remarks" binding:"required"`
}

type TransferRequest struct {
	TargetUser string  `json:"target_user" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Remarks    string  `json:"remarks" binding:"required"`
}

type TopUpResponse struct {
	Status string      `json:"status"`
	Result TopUpResult `json:"result"`
}

type PaymentResponse struct {
	Status string        `json:"status"`
	Result PaymentResult `json:"result"`
}

type TransferResponse struct {
	Status string         `json:"status"`
	Result TransferResult `json:"result"`
}

type GetTransactionsResponse struct {
	Status string                `json:"status"`
	Result []TransactionResponse `json:"result"`
}

type TopUpResult struct {
	TopUpID       string    `json:"top_up_id"`
	AmountTopUp   float64   `json:"amount_top_up"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedAt     time.Time `json:"created_date"`
}

type PaymentResult struct {
	PaymentID     string    `json:"payment_id"`
	Amount        float64   `json:"amount"`
	Remarks       string    `json:"remarks"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedAt     time.Time `json:"created_date"`
}

type TransferResult struct {
	TransferID    string    `json:"transfer_id"`
	Amount        float64   `json:"amount"`
	Remarks       string    `json:"remarks"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedAt     time.Time `json:"created_date"`
}

type TransactionResponse struct {
	TransferID      string    `json:"transfer_id"`
	Status          string    `json:"status"`
	UserID          string    `json:"user_id"`
	TransactionType string    `json:"transaction_type"`
	Amount          float64   `json:"amount"`
	Remarks         string    `json:"remarks"`
	BalanceBefore   float64   `json:"balance_before"`
	BalanceAfter    float64   `json:"balance_after"`
	CreatedAt       time.Time `json:"created_date"`
}
