package services

import (
	"testing"

	"github.com/ycChu711/receipt-processor/models"
)

func TestCalculatePoints(t *testing.T) {

	const (
		defaultTestDate = "2022-01-01"
		defaultTestTime = "13:01"
	)

	tests := []struct {
		name     string
		receipt  models.Receipt
		expected int64
	}{
		{
			name: "Target example",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
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
			name: "M&M example",
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
			name: "Retailer name points",
			receipt: models.Receipt{
				Retailer:     "A&W",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 83,
		},
		{
			name: "Store with special chars",
			receipt: models.Receipt{
				Retailer:     "& - &",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 81,
		},
		{
			name: "Mixed case name",
			receipt: models.Receipt{
				Retailer:     "MiXeD cAsE sToRe",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 95,
		},
		{
			name: "Long store name",
			receipt: models.Receipt{
				Retailer:     "asdfghjklzxcvbnm1234567890zqwertyuioplkjhgfdsazxcvbnmklhy-111122223333444455556666",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 162,
		},
		{
			name: "Round dollar amount",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items:        []models.Item{{ShortDescription: "Item", Price: "5.00"}},
				Total:        "5.00",
			},
			expected: 85,
		},
		{
			name: "Multiple of 0.25",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items:        []models.Item{{ShortDescription: "Item", Price: "5.25"}},
				Total:        "5.25",
			},
			expected: 35,
		},
		{
			name: "Free item",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items:        []models.Item{{ShortDescription: "Free Item", Price: "0.00"}},
				Total:        "0.00",
			},
			expected: 85,
		},
		{
			name: "0.25",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items:        []models.Item{{ShortDescription: "Tiny", Price: "0.25"}},
				Total:        "0.25",
			},
			expected: 35,
		},
		{
			name: "Total ending in .50",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items:        []models.Item{{ShortDescription: "Half Dollar", Price: "1.50"}},
				Total:        "1.50",
			},
			expected: 35,
		},
		{
			name: "Total ending in .75",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items:        []models.Item{{ShortDescription: "Three Quarters", Price: "1.75"}},
				Total:        "1.75",
			},
			expected: 35,
		},
		{
			name: "Not multiple of 0.25",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items:        []models.Item{{ShortDescription: "Odd Price", Price: "1.37"}},
				Total:        "1.37",
			},
			expected: 11,
		},
		{
			name: "pairs of items",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
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
		},
		{
			name: "One item",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items:        []models.Item{{ShortDescription: "Single Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 85,
		},
		{
			name: "3 items",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items: []models.Item{
					{ShortDescription: "Item1", Price: "1.00"},
					{ShortDescription: "Item2", Price: "1.00"},
					{ShortDescription: "Item3", Price: "1.00"},
				},
				Total: "3.00",
			},
			expected: 90,
		},
		{
			name: "10 items",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items: []models.Item{
					{ShortDescription: "Item1", Price: "1.00"},
					{ShortDescription: "Item2", Price: "1.00"},
					{ShortDescription: "Item3", Price: "1.00"},
					{ShortDescription: "Item4", Price: "1.00"},
					{ShortDescription: "Item5", Price: "1.00"},
					{ShortDescription: "Item6", Price: "1.00"},
					{ShortDescription: "Item7", Price: "1.00"},
					{ShortDescription: "Item8", Price: "1.00"},
					{ShortDescription: "Item9", Price: "1.00"},
					{ShortDescription: "Item10", Price: "1.00"},
				},
				Total: "10.00",
			},
			expected: 111,
		},
		{
			name: "Description length multiple of 3",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items: []models.Item{
					{ShortDescription: "123456789", Price: "3.00"},
				},
				Total: "3.00",
			},
			expected: 86,
		},
		{
			name: "Spaces need trimming",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items:        []models.Item{{ShortDescription: "   ABC ", Price: "5.00"}}, // Trimmed length: 3
				Total:        "5.00",
			},
			expected: 86,
		},
		{
			name: "Multiple descriptions with len multiple of 3",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items: []models.Item{
					{ShortDescription: "ABC", Price: "5.00"},
					{ShortDescription: "DEFGHI", Price: "10.00"},
					{ShortDescription: "JKL", Price: "15.00"},
				},
				Total: "30.00",
			},
			expected: 96,
		},
		{
			name: "Odd purchase day",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: defaultTestTime,
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 85,
		},
		{
			name: "1-31",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-31",
				PurchaseTime: defaultTestTime,
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 85,
		},
		{
			name: "30th day",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-04-30",
				PurchaseTime: defaultTestTime,
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 79,
		},
		{
			name: "Purchase between 2 and 4 pm",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: "14:30",
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 95,
		},
		{
			name: "2pm",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: "14:00",
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 85,
		},
		{
			name: "just before 2pm - 13:59",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: "13:59",
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 85,
		},
		{
			name: "just before 4pm - 15:59",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: "15:59",
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 95,
		},
		{
			name: "4pm",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: defaultTestDate,
				PurchaseTime: "16:00",
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 85,
		},
		{
			name: "low points senario",
			receipt: models.Receipt{
				Retailer:     "A",
				PurchaseDate: "2022-02-02",
				PurchaseTime: "12:00",
				Items:        []models.Item{{ShortDescription: "Item", Price: "0.37"}},
				Total:        "0.37",
			},
			expected: 1,
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
