package services

import (
	"github.com/google/uuid"
	"github.com/ycChu711/receipt-processor/models"
	"github.com/ycChu711/receipt-processor/repository"
)

// manages receipt processing and point calculations
type ReceiptService struct {
	storage repository.ReceiptStorage
}

// create new service with given storage
func NewReceiptService(storage repository.ReceiptStorage) *ReceiptService {
	return &ReceiptService{
		storage: storage,
	}
}

// Processes a receipt and returns the ID
// generate unique id -> calculate points -> save receipt and points -> return id
func (s *ReceiptService) ProcessReceipt(receipt models.Receipt) (string, error) {

	id := uuid.New().String()

	points := CalculatePoints(&receipt)

	err := s.storage.SaveReceipt(id, receipt, points)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *ReceiptService) GetPoints(id string) (int64, bool) {
	return s.storage.GetPoints(id)
}
