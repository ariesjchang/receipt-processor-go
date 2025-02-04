package services

import (
	"testing"

	"github.com/ariesjchang/receipt-processor-go/internal/models"
)

func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		name     string
		receipt  models.Receipt
		expected int
	}{
		{
			name: "Basic Receipt - Various Rules",
			receipt: models.Receipt{
				Retailer:     "Target", // 6 alphanumeric chars → 6 points
				Total:        "35.00",  // Round total → +50 points, Multiple of 0.25 → +25 points
				PurchaseDate: "2022-01-01", // Odd day → +6 points
				PurchaseTime: "14:30", // Between 2-4 PM → +10 points
				Items: []models.Item{ // 1 pair of items → +5 points
					{ShortDescription: "Item A", Price: "12.00"}, // Length 6 (multiple of 3) → +3 points
					{ShortDescription: "Item B", Price: "3.00"},  // Length 6 (multiple of 3) → +1 point
					{ShortDescription: "Long Item Desc", Price: "4.75"}, // Length 14
				},
			},
			expected: 106, // 6 + 50 + 25 + 6 + 10 + 5 + 3 + 1
		},
		{
			name: "Only Retailer Points",
			receipt: models.Receipt{
				Retailer: "BestBuy123",
				Total:    "19.99",
				Items:    []models.Item{},
			},
			expected: 10, // 10 alphanumeric characters
		},
		{
			name: "Round Dollar & Multiple of 0.25",
			receipt: models.Receipt{
				Total: "20.00",
			},
			expected: 50 + 25, // 50 (round) + 25 (0.25 multiple)
		},
		{
			name: "Not Round Dollar & Multiple of 0.25",
			receipt: models.Receipt{
				Total: "20.25",
			},
			expected: 25, // 0.25 multiple
		},
		{
			name: "Not Round Dollar & Not Multiple of 0.25",
			receipt: models.Receipt{
				Total: "20.45",
			},
			expected: 0,
		},
		{
			name: "Five Items - Two Pairs Bonus",
			receipt: models.Receipt{
				Items: []models.Item{
					{ShortDescription: "Item1", Price: "5.50"},
					{ShortDescription: "Item2", Price: "2.25"},
					{ShortDescription: "Item3", Price: "4.75"},
					{ShortDescription: "Item4", Price: "1.05"},
					{ShortDescription: "Item5", Price: "2.25"},
				},
			},
			expected: 10, // 2 pairs of items → 10 points
		},
		{
			name: "Two Items - Pair Bonus",
			receipt: models.Receipt{
				Items: []models.Item{
					{ShortDescription: "Item1", Price: "5.50"},
					{ShortDescription: "Item2", Price: "2.25"},
				},
			},
			expected: 5, // 1 pair of items → 5 points
		},
		{
			name: "One Item - No Pair Bonus",
			receipt: models.Receipt{
				Items: []models.Item{
					{ShortDescription: "Item1", Price: "5.50"},
				},
			},
			expected: 0,
		},
		{
			name: "Item Description Length Multiple of 3",
			receipt: models.Receipt{
				Items: []models.Item{
					{ShortDescription: "BBBBBB", Price: "20.00"}, // Length 6, price * 0.2 = 4 points
				},
			},
			expected: 4, // 6 points total
		},
		{
			name: "Item Description Length Not Multiple of 3",
			receipt: models.Receipt{
				Items: []models.Item{
					{ShortDescription: "BBB BBB", Price: "20.00"}, // Length 7
				},
			},
			expected: 0,
		},
		{
			name: "Odd Day Bonus",
			receipt: models.Receipt{
				PurchaseDate: "2022-03-15",
			},
			expected: 6, // Odd day
		},
		{
			name: "Even Day No Bonus",
			receipt: models.Receipt{
				PurchaseDate: "2022-03-14",
			},
			expected: 0,
		},
		{
			name: "Time Bonus - Within 2:00PM-4:00PM",
			receipt: models.Receipt{
				PurchaseTime: "14:01",
			},
			expected: 10, // Within the time range
		},
		{
			name: "No Time Bonus at 2:00PM",
			receipt: models.Receipt{
				PurchaseTime: "14:00",
			},
			expected: 0, // Not after 2:00PM
		},
		{
			name: "No Time Bonus at 4:00PM",
			receipt: models.Receipt{
				PurchaseTime: "16:00",
			},
			expected: 0, // Not before 4:00PM
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculatePoints(tt.receipt)
			if got != tt.expected {
				t.Errorf("CalculatePoints() = %d; expected %d", got, tt.expected)
			}
		})
	}
}