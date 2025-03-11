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
			name: "Edge Case - Retailer Name with Only Valid Special Characters",
			receipt: models.Receipt{
				Retailer:     "& - &",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 81,
			// 0 (no alphanumeric in retailer) + 50 (round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 0 (time)
		},
		{
			name: "Edge Case - Mixed Case Retailer Name",
			receipt: models.Receipt{
				Retailer:     "MiXeD cAsE sToRe",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 95,
			// 14 (alphanumeric chars) + 50 (round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 0 (time)
		},
		{
			name: "Edge Case - Very Long Retailer Name",
			receipt: models.Receipt{
				Retailer:     "Super-Extremely-Long-Retailer-Name-With-Many-Characters-For-Testing-Purposes-12345",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 152,
			// 71 (alphanumeric chars only) + 50 (round dollar) + 25 (multiple of 0.25)
			// + 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 0 (time)
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
			name: "Edge Case - Zero Total",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Free Item", Price: "0.00"}},
				Total:        "0.00",
			},
			expected: 85,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25)
			// + 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 0 (time)
		},
		{
			name: "Edge Case - Very Small Total Multiple of 0.25",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Tiny", Price: "0.25"}},
				Total:        "0.25",
			},
			expected: 35,
			// 4 (retailer) + 0 (not round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 0 (time)
		},
		{
			name: "Edge Case - Total Ending in .50",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Half Dollar", Price: "1.50"}},
				Total:        "1.50",
			},
			expected: 35,
			// 4 (retailer) + 0 (not round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 0 (time)
		},
		{
			name: "Edge Case - Total Ending in .75",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Three Quarters", Price: "1.75"}},
				Total:        "1.75",
			},
			expected: 35,
			// 4 (retailer) + 0 (not round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 0 (time)
		},
		{
			name: "Edge Case - Total Not Multiple of 0.25",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Odd Price", Price: "1.37"}},
				Total:        "1.37",
			},
			expected: 11,
			// 4 (retailer) + 0 (not round dollar) + 0 (not multiple of 0.25) +
			// 0 (0 item pairs) + 1 (desc points) + 6 (odd day) + 0 (time)
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
			name: "Edge Case - Single Item",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Single Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 85,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 0 (time)
		},
		{
			name: "Edge Case - Odd Number of Items",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "Item1", Price: "1.00"},
					{ShortDescription: "Item2", Price: "1.00"},
					{ShortDescription: "Item3", Price: "1.00"},
				},
				Total: "3.00",
			},
			expected: 90,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) + 5 (1 item pair) +
			// 0 (desc points) + 6 (odd day) + 0 (time)
		},
		{
			name: "Edge Case - Many Items",
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
					{ShortDescription: "Item6", Price: "1.00"},
					{ShortDescription: "Item7", Price: "1.00"},
					{ShortDescription: "Item8", Price: "1.00"},
					{ShortDescription: "Item9", Price: "1.00"},
					{ShortDescription: "Item10", Price: "1.00"},
				},
				Total: "10.00",
			},
			expected: 111,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) + 25 (5 item pairs) +
			// 1 (desc points) + 6 (odd day) + 0 (time)
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
			name: "Edge Case - Description with Spaces Needing Trim",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "   ABC ", Price: "5.00"}}, // Trimmed length: 3
				Total:        "5.00",
			},
			expected: 86,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 1 (5.00*0.2=1) + 6 (odd day) + 0 (time)
		},
		{
			name: "Edge Case - Multiple Descriptions of Length Multiple of 3",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "ABC", Price: "5.00"},
					{ShortDescription: "DEFGHI", Price: "10.00"},
					{ShortDescription: "JKL", Price: "15.00"},
				},
				Total: "30.00",
			},
			expected: 96,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) + 5 (1 item pair) +
			// 6 (desc points: 1+2+3) + 6 (odd day) + 0 (time)
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
			name: "Edge Case - Last Day of Month",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-31",
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 85,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 0 (time)
		},
		{
			name: "Edge Case - 30th Day",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-04-30", // 30th is even
				PurchaseTime: "13:01",
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 79,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc points) + 0 (even day) + 0 (time)
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
		{
			name: "Edge Case - Exactly 2:00 PM",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "14:00", // 2:00 PM
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 95,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 10 (time is 2-4 PM)
		},
		{
			name: "Edge Case - Just Before 2:00 PM",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:59", // 1:59 PM
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 85,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 0 (time not 2-4 PM)
		},
		{
			name: "Edge Case - Just Before 4:00 PM",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "15:59", // 3:59 PM
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 95,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) + 0 (0 item pairs) +
			// 0 (desc points) + 6 (odd day) + 10 (time is 2-4 PM)
		},
		{
			name: "Edge Case - Exactly 4:00 PM",
			receipt: models.Receipt{
				Retailer:     "Shop",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "16:00", // 4:00 PM
				Items:        []models.Item{{ShortDescription: "Item", Price: "1.00"}},
				Total:        "1.00",
			},
			expected: 85,
			// 4 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc points) + 6 (odd day) + 0 (time not 2-4 PM)
		},
		{
			name: "Edge Case - Maximum Points Scenario",
			receipt: models.Receipt{
				Retailer:     "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890",
				PurchaseDate: "2023-01-31",
				PurchaseTime: "14:30",
				Items: []models.Item{
					{ShortDescription: "abc", Price: "10.00"},
					{ShortDescription: "def", Price: "10.00"},
					{ShortDescription: "ghi", Price: "10.00"},
					{ShortDescription: "jkl", Price: "10.00"},
					{ShortDescription: "mno", Price: "10.00"},
					{ShortDescription: "pqr", Price: "10.00"},
				},
				Total: "60.00",
			},
			expected: 154,
			// 36 (retailer) + 50 (round dollar) + 25 (multiple of 0.25) + 15 (3 item pairs) +
			// 12 (desc points: 6*2=12) + 6 (odd day) + 10 (time)
		},
		{
			name: "Edge Case - Minimum Valid Points Scenario",
			receipt: models.Receipt{
				Retailer:     "A",
				PurchaseDate: "2022-02-02",
				PurchaseTime: "12:00",
				Items:        []models.Item{{ShortDescription: "Item", Price: "0.37"}},
				Total:        "0.37",
			},
			expected: 1,
			// 1 (retailer) + 0 (not round dollar) + 0 (not multiple of 0.25) +
			// 0 (0 item pairs) + 0 (desc not multiple of 3) + 0 (even day) + 0 (time)
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
