package services

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ycChu711/receipt-processor/models"
)

// CalculatePoints calculates the points for a receipt
func CalculatePoints(receipt *models.Receipt) int64 {
	var points int64 = 0

	fmt.Printf("\n========== POINTS CALCULATION FOR %s ==========\n", receipt.Retailer)

	// Rule 1: One point for every alphanumeric character in the retailer name
	retailerName := receipt.Retailer
	var alphaNumCount int64 = 0
	alphaNumChars := ""
	for _, char := range retailerName {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			alphaNumCount++
			alphaNumChars += string(char)
		}
	}
	points += alphaNumCount

	fmt.Printf("Rule 1 - Retailer name: '%s'\n", receipt.Retailer)
	fmt.Printf("        Alphanumeric chars: '%s' (count: %d)\n", alphaNumChars, alphaNumCount)
	fmt.Printf("        Points from Rule 1: %d\n", alphaNumCount)
	fmt.Printf("        Running total: %d\n\n", points)

	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	totalFloat, _ := strconv.ParseFloat(receipt.Total, 64)
	if totalFloat == float64(int64(totalFloat)) {
		points += 50

		fmt.Printf("Rule 2 - Total $%s is a round dollar amount\n", receipt.Total)
		fmt.Printf("        Points from Rule 2: 50\n")
	} else {
		fmt.Printf("Rule 2 - Total $%s is not a round dollar amount\n", receipt.Total)
		fmt.Printf("        Points from Rule 2: 0\n")
	}
	fmt.Printf("        Running total: %d\n\n", points)

	// Rule 3: 25 points if the total is a multiple of 0.25.
	if math.Mod(totalFloat*100, 25) == 0 {
		points += 25

		fmt.Printf("Rule 3 - Total $%s is a multiple of 0.25\n", receipt.Total)
		fmt.Printf("        Points from Rule 3: 25\n")
	} else {
		fmt.Printf("Rule 3 - Total $%s is NOT a multiple of 0.25\n", receipt.Total)
		fmt.Printf("        Points from Rule 3: 0\n")
	}
	fmt.Printf("        Running total: %d\n\n", points)

	// Rule 4: 5 points for every two items on the receipt.
	itemPairs := len(receipt.Items) / 2
	pairPoints := int64(itemPairs * 5)
	points += pairPoints
	fmt.Printf("Rule 4 - %d items = %d pairs\n", len(receipt.Items), itemPairs)
	fmt.Printf("        Points from Rule 4: %d\n", pairPoints)
	fmt.Printf("        Running total: %d\n\n", points)

	// Rule 5: If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	fmt.Printf("Rule 5 - Items with description length multiple of 3:\n")
	var descPoints int64 = 0
	for i, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		trimmedLen := len(trimmedDesc)
		fmt.Printf("        Item %d: '%s' (trimmed length: %d)\n",
			i+1, trimmedDesc, trimmedLen)

		if trimmedLen > 0 && trimmedLen%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			calculation := price * 0.2
			itemPoints := int64(math.Ceil(calculation))
			descPoints += itemPoints
			fmt.Printf("        Multiple of 3! Price: $%s, Calculation: %.2f * 0.2 = %.2f, Rounded up: %d\n",
				item.Price, price, calculation, itemPoints)
		} else {
			fmt.Printf("        Not a multiple of 3, no points\n")
		}
	}
	points += descPoints
	fmt.Printf("        Total points from Rule 5: %d\n", descPoints)
	fmt.Printf("        Running total: %d\n\n", points)

	// Rule 7: 6 points if the day in the purchase date is odd.
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	day := purchaseDate.Day()
	if day%2 == 1 {
		points += 6
		fmt.Printf("Rule 7 - Purchase day %d is odd\n", day)
		fmt.Printf("        Points from Rule 7: 6\n")
	} else {
		fmt.Printf("Rule 7 - Purchase day %d is even\n", day)
		fmt.Printf("        Points from Rule 7: 0\n")
	}
	fmt.Printf("        Running total: %d\n\n", points)

	// Rule 8: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	hour := purchaseTime.Hour()
	minute := purchaseTime.Minute()

	// Check if time is between 2:00 PM (14:00) and 4:00 PM (16:00)
	inTimeRange := (hour == 14 && minute >= 0) || (hour == 15)

	if inTimeRange {
		points += 10
		fmt.Printf("Rule 8 - Purchase time %s is between 2:00 PM and 4:00 PM\n", receipt.PurchaseTime)
		fmt.Printf("        Points from Rule 8: 10\n")
	} else {
		fmt.Printf("Rule 8 - Purchase time %s is NOT between 2:00 PM and 4:00 PM\n", receipt.PurchaseTime)
		fmt.Printf("        Points from Rule 8: 0\n")
	}
	fmt.Printf("        Running total: %d\n\n", points)

	fmt.Printf("FINAL TOTAL: %d points\n", points)
	fmt.Printf("=============================================\n")

	return points
}
