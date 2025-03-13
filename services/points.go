package services

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ycChu711/receipt-processor/models"
	"github.com/ycChu711/receipt-processor/utils"
)

func CalculatePoints(receipt *models.Receipt) int64 {
	utils.Logger.WithField("retailer", receipt.Retailer).Info("Starting points calculation")

	retailerPoints := calculateRetailerNamePoints(receipt.Retailer)
	roundDollarPoints := calculateRoundDollarPoints(receipt.Total)
	quarterMultiplePoints := calculateQuarterMultiplePoints(receipt.Total)
	itemPairPoints := calculateItemPairPoints(len(receipt.Items))
	descriptionPoints := calculateDescriptionLengthPoints(receipt.Items)
	oddDayPoints := calculateOddDayPoints(receipt.PurchaseDate)
	timeRangePoints := calculateTimeRangePoints(receipt.PurchaseTime)

	totalPoints := retailerPoints + roundDollarPoints + quarterMultiplePoints +
		itemPairPoints + descriptionPoints + oddDayPoints + timeRangePoints

	utils.Logger.WithFields(logrus.Fields{
		"retailer":     receipt.Retailer,
		"final_points": totalPoints,
	}).Info("Completed points calculation")

	return totalPoints
}

// Rule 1: One point for every alphanumeric character in the retailer name
func calculateRetailerNamePoints(name string) int64 {
	points := int64(0)
	validChars := ""

	for _, char := range name {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			points++
			validChars += string(char)
		}
	}

	utils.Logger.WithFields(logrus.Fields{
		"rule":             "1",
		"retailer":         name,
		"alphanumeric":     validChars,
		"count":            points,
		"points_from_rule": points,
	}).Debug("Retailer name points rule")

	return points
}

// Rule 2: 50 points if the total is a round dollar amount with no cents
func calculateRoundDollarPoints(total string) int64 {
	totalFloat, _ := strconv.ParseFloat(total, 64)
	isRoundDollar := totalFloat == float64(int64(totalFloat))
	points := int64(0)

	if isRoundDollar {
		points = 50
		utils.Logger.Debugf("Rule 2: %s is a round dollar amount", total)
	} else {
		utils.Logger.Debugf("Rule 2: %s is not a round dollar amount", total)
	}
	return points
}

// Rule 3: 25 points if the total is a multiple of 0.25
func calculateQuarterMultiplePoints(total string) int64 {
	totalFloat, _ := strconv.ParseFloat(total, 64)
	isMultipleOfQuarter := math.Mod(totalFloat*100, 25) == 0
	points := int64(0)

	if isMultipleOfQuarter {
		points = 25
		utils.Logger.Debug("Total is a multiple of 0.25")
	}

	return points
}

// Rule 4: 5 points for every two items on the receipt
func calculateItemPairPoints(itemCount int) int64 {
	pairs := itemCount / 2
	points := int64(pairs * 5)

	utils.Logger.Debugf("%d items on the receipt, %d pairs, %d points", itemCount, pairs, points)
	return points
}

// Rule 5: If the trimmed length of the item description is a multiple of 3,
// multiply the price by 0.2 and round up to the nearest integer
func calculateDescriptionLengthPoints(items []models.Item) int64 {
	var totalPoints int64 = 0

	for i, item := range items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		trimmedLen := len(trimmedDesc)
		itemPoints := int64(0)

		if trimmedLen > 0 && trimmedLen%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			itemPoints = int64(math.Ceil(price * 0.2))
			totalPoints += itemPoints

		} else {
			utils.Logger.Debugf("Item %d description length is not a multiple of 3", i)
		}
	}

	utils.Logger.Debugf("Total points from rule 5: %d", totalPoints)

	return totalPoints
}

// Rule 6: 6 points if the day in the purchase date is odd
func calculateOddDayPoints(purchaseDate string) int64 {
	date, _ := time.Parse("2006-01-02", purchaseDate)
	day := date.Day()
	isOddDay := day%2 == 1
	points := int64(0)

	if isOddDay {
		points = 6
	}

	return points
}

// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm
func calculateTimeRangePoints(purchaseTime string) int64 {
	time, _ := time.Parse("15:04", purchaseTime)
	hour := time.Hour()
	minute := time.Minute()

	inTimeRange := (hour == 14 && minute > 0) || (hour == 15)
	points := int64(0)

	if inTimeRange {
		points = 10
	}

	return points
}
