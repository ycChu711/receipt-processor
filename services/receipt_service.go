package services

import (
	"github.com/google/uuid"
	"github.com/ycChu711/receipt-processor/models"
	"github.com/ycChu711/receipt-processor/repository"
)

// ReceiptService handles receipt processing logic
type ReceiptService struct {
	storage repository.ReceiptStorage
}

// NewReceiptService creates a new receipt service
func NewReceiptService(storage repository.ReceiptStorage) *ReceiptService {
	return &ReceiptService{
		storage: storage,
	}
}

// ProcessReceipt processes a receipt and returns the ID
func (s *ReceiptService) ProcessReceipt(receipt models.Receipt) (string, error) {
	// Generate a unique ID
	id := uuid.New().String()

	// Calculate points
	points := CalculatePoints(&receipt)

	// Save receipt and points
	err := s.storage.SaveReceipt(id, receipt, points)
	if err != nil {
		return "", err
	}

	return id, nil
}

// GetPoints gets the points for a receipt by ID
func (s *ReceiptService) GetPoints(id string) (int64, bool) {
	return s.storage.GetPoints(id)
}
