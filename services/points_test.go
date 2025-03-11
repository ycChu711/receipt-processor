package services

import (
	"testing"

	"github.com/ycChu711/receipt-processor/models"
)

func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		name     string
		receipt  models.Receipt
		expected int64
	}{
		{
			name: "Target Example",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
					{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
					{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
					{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
					{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
				},
				Total: "35.35",
			},
			expected: 28,
		}, {
			name: "M&M Corner Market Example",
			receipt: models.Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: "2022-03-20",
				PurchaseTime: "14:33",
				Items: []models.Item{
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
				},
				Total: "9.00",
			},
			expected: 109,
		},
		{
			name: "Testing Rule 1: Retailer name points",
			receipt: models.Receipt{
				Retailer:     "A&W",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 83,
			// 2 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 0 (time)
		},
		{
			name: "Testing Rule 2: Round dollar amount",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Item", Price: "5.00"}},
				Total:        "5.00",
			},
			expected: 85,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 0 (time)
		},
		{
			name: "Testing Rule 3: Multiple of 0.25",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Item", Price: "5.25"}},
				Total:        "5.25",
			},
			expected: 35,
			// 4 (retailer) + 0 (not round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 0 (time)
		},
		{
			name: "Testing Rule 4: Points for pairs of items",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "Item1", Price: "1.00"},
					{ShortDescription: "Item2", Price: "1.00"},
					{ShortDescription: "Item3", Price: "1.00"},
					{ShortDescription: "Item4", Price: "1.00"},
					{ShortDescription: "Item5", Price: "1.00"},
				},
				Total: "5.00",
			},
			expected: 95,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) +
			// 10 (2 item pairs) + 0 (desc points) + 6 (odd day) + 0 (time)
		},
		{
			name: "Testing Rule 5: Description length multiple of 3",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "123456789", Price: "3.00"},
				},
				Total: "3.00",
			},
			expected: 86,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 1 (desc points) + 6 (odd day) + 0 (time)
		},
		{
			name: "Testing Rule 7: Odd purchase day",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 85,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 0 (time)
		},
		{
			name: "Testing Rule 8: Purchase time between 2:00PM and 4:00PM",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "14:30", // 2:30PM
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 95,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 10 (time)
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			points := CalculatePoints(&tc.receipt)
			if points != tc.expected {
				t.Errorf("Expected %d points, got %d", tc.expected, points)
			}
		})
	}
}
