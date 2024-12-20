package worker

import (
	"casion/internal/models"
	"log"
	"time"

	"gorm.io/gorm"
)

func ProcessTransfers(db *gorm.DB) {
	for {
		var pendingTransfers []models.Transaction
		result := db.Where("type = ? AND status = ?", "transfer", "pending").Find(&pendingTransfers)
		if result.Error != nil {
			log.Printf("Error fetching pending transfers: %v", result.Error)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, transfer := range pendingTransfers {
			// Start transaction
			tx := db.Begin()

			// Get sender
			var sender models.User
			if err := tx.First(&sender, "user_id = ?", transfer.UserID).Error; err != nil {
				tx.Rollback()
				transfer.Status = "failed"
				transfer.Remarks = "Sender not found"
				db.Save(&transfer)
				continue
			}

			// Get receiver
			var receiver models.User
			if err := tx.First(&receiver, "user_id = ?", transfer.TargetUserID).Error; err != nil {
				tx.Rollback()
				transfer.Status = "failed"
				transfer.Remarks = "Receiver not found"
				db.Save(&transfer)
				continue
			}

			// Check balance
			if sender.Balance < transfer.Amount {
				tx.Rollback()
				transfer.Status = "failed"
				transfer.Remarks = "Insufficient balance"
				db.Save(&transfer)
				continue
			}

			// Update balances
			sender.Balance -= transfer.Amount
			receiver.Balance += transfer.Amount

			// Save changes
			if err := tx.Save(&sender).Error; err != nil {
				tx.Rollback()
				transfer.Status = "failed"
				transfer.Remarks = "Failed to update sender balance"
				db.Save(&transfer)
				continue
			}

			if err := tx.Save(&receiver).Error; err != nil {
				tx.Rollback()
				transfer.Status = "failed"
				transfer.Remarks = "Failed to update receiver balance"
				db.Save(&transfer)
				continue
			}

			// Update transfer status
			transfer.Status = "success"
			if err := tx.Save(&transfer).Error; err != nil {
				tx.Rollback()
				transfer.Status = "failed"
				transfer.Remarks = "Failed to update transfer status"
				db.Save(&transfer)
				continue
			}

			tx.Commit()
		}

		time.Sleep(5 * time.Second)
	}
}
