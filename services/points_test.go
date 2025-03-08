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
			expected: 53, // 3 alphanumeric chars (A, W) + 50 points for round dollar
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
			expected: 54, // 4 points for retailer + 50 points for round dollar amount
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
			expected: 29, // 4 points for retailer + 25 points for multiple of 0.25
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
			expected: 64, // 4 points for retailer + 50 points for round dollar + 10 points for 2 pairs
		},
		{
			name: "Testing Rule 5: Description length multiple of 3",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "123456789", Price: "3.00"}, // Length is 9, a multiple of 3
				},
				Total: "3.00",
			},
			expected: 55, // 4 points for retailer + 50 points for round dollar + 1 point (3.00 * 0.2 = 0.6, rounded up to 1)
		},
		{
			name: "Testing Rule 6: Total greater than 10.00",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Item", Price: "12.00"}},
				Total:        "12.00",
			},
			expected: 59, // 4 points for retailer + 50 points for round dollar + 5 points for total > 10.00
		},
		{
			name: "Testing Rule 7: Odd purchase day",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01", // 1st is odd
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 60, // 4 points for retailer + 50 points for round dollar + 6 points for odd day
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
			expected: 70, // 4 points for retailer + 50 points for round dollar + 6 points for odd day + 10 points for time
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
