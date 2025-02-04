package services

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ariesjchang/receipt-processor-go/internal/models"
)

// CalculatePoints computes the points earned from a receipt based on predefined rules.
func CalculatePoints(receipt models.Receipt) int {
	points := 0

	// 1. One point for each alphanumeric character in retailer name
	points += countAlphanumeric(receipt.Retailer)

	// 2. Check if the total is a round number (50 points) and/or multiple of 0.25 (25 points)
	if total, err := strconv.ParseFloat(receipt.Total, 64); err == nil {
		if math.Mod(total, 1.0) == 0 {
			points += 50
		}
		if int(total*100)%25 == 0 {
			points += 25
		}
	}

	// 3. 5 points for every two items
	points += (len(receipt.Items) / 2) * 5

	// 4. Extra points for item description length being a multiple of 3
	for _, item := range receipt.Items {
		desc := strings.TrimSpace(item.ShortDescription)
		if len(desc)%3 == 0 {
			if price, err := strconv.ParseFloat(item.Price, 64); err == nil {
				points += int(math.Ceil(price * 0.2))
			}
		}
	}

	// 5. 6 points if the purchase day is odd
	if dateParts := strings.Split(receipt.PurchaseDate, "-"); len(dateParts) == 3 {
		if day, err := strconv.Atoi(dateParts[2]); err == nil && day%2 != 0 {
			points += 6
		}
	}

	// 6. 10 points if purchase time is between 2:00pm and 4:00pm
	if t, err := time.Parse("15:04", receipt.PurchaseTime); err == nil {
		if t.Hour() >= 14 && t.Hour() < 16 {
			points += 10
		}
	}

	return points
}

// countAlphanumeric counts only alphanumeric characters in a string.
func countAlphanumeric(s string) int {
	count := 0
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			count++
		}
	}
	return count
}
