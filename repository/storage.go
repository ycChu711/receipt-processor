package repository

import (
	"sync"

	"github.com/ycChu711/receipt-processor/models"
)

type ReceiptStorage interface {
	SaveReceipt(id string, receipt models.Receipt, points int64) error
	GetPoints(id string) (int64, bool)
}

type InMemoryStorage struct {
	receiptsWithPoints map[string]models.ReceiptWithPoints
	mutex              *sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		receiptsWithPoints: map[string]models.ReceiptWithPoints{},
		mutex:              &sync.RWMutex{},
	}
}

func (s *InMemoryStorage) SaveReceipt(id string, receipt models.Receipt, points int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.receiptsWithPoints[id] = models.ReceiptWithPoints{
		Receipt: receipt,
		Points:  points,
	}
	return nil
}

func (s *InMemoryStorage) GetReceipt(id string) (models.Receipt, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	receiptWithPoints, exists := s.receiptsWithPoints[id]
	if !exists {
		return models.Receipt{}, false
	}
	return receiptWithPoints.Receipt, true
}

func (s *InMemoryStorage) GetPoints(id string) (int64, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	receiptWithPoints, exists := s.receiptsWithPoints[id]
	if !exists {
		return 0, false
	}
	return receiptWithPoints.Points, true
}
