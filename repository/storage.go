package repository

import (
	"sync"

	"github.com/ycChu711/receipt-processor/models"
)

// ReceiptStorage defines operations for persisting receipts
type ReceiptStorage interface {
	SaveReceipt(id string, receipt models.Receipt, points int64) error
	GetPoints(id string) (int64, bool)
}

type InMemoryStorage struct {
	receipts map[string]models.Receipt
	points   map[string]int64
	mutex    *sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		receipts: make(map[string]models.Receipt),
		points:   make(map[string]int64),
		mutex:    &sync.RWMutex{},
	}
}

func (s *InMemoryStorage) SaveReceipt(id string, receipt models.Receipt, points int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.receipts[id] = receipt
	s.points[id] = points
	return nil
}

func (s *InMemoryStorage) GetPoints(id string) (int64, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	points, exists := s.points[id]
	return points, exists
}
