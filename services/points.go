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
func calculateRetailerNamePoints(retailerName string) int64 {
	var alphaNumCount int64 = 0
	alphaNumChars := ""

	for _, char := range retailerName {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			alphaNumCount++
			alphaNumChars += string(char)
		}
	}

	utils.Logger.WithFields(logrus.Fields{
		"rule":             "1",
		"retailer":         retailerName,
		"alphanumeric":     alphaNumChars,
		"count":            alphaNumCount,
		"points_from_rule": alphaNumCount,
	}).Debug("Applied retailer name points rule")

	return alphaNumCount
}

// Rule 2: 50 points if the total is a round dollar amount with no cents
func calculateRoundDollarPoints(total string) int64 {
	totalFloat, _ := strconv.ParseFloat(total, 64)
	isRoundDollar := totalFloat == float64(int64(totalFloat))
	points := int64(0)

	if isRoundDollar {
		points = 50
	}

	utils.Logger.WithFields(logrus.Fields{
		"rule":             "2",
		"total":            total,
		"is_round_dollar":  isRoundDollar,
		"points_from_rule": points,
	}).Debug("Applied round dollar amount rule")

	return points
}

// Rule 3: 25 points if the total is a multiple of 0.25
func calculateQuarterMultiplePoints(total string) int64 {
	totalFloat, _ := strconv.ParseFloat(total, 64)
	isMultipleOfQuarter := math.Mod(totalFloat*100, 25) == 0
	points := int64(0)

	if isMultipleOfQuarter {
		points = 25
	}

	utils.Logger.WithFields(logrus.Fields{
		"rule":               "3",
		"total":              total,
		"is_multiple_of_025": isMultipleOfQuarter,
		"points_from_rule":   points,
	}).Debug("Applied multiple of 0.25 rule")

	return points
}

// Rule 4: 5 points for every two items on the receipt
func calculateItemPairPoints(itemCount int) int64 {
	itemPairs := itemCount / 2
	points := int64(itemPairs * 5)

	utils.Logger.WithFields(logrus.Fields{
		"rule":             "4",
		"items_count":      itemCount,
		"pairs":            itemPairs,
		"points_from_rule": points,
	}).Debug("Applied item pairs rule")

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

		itemLogger := utils.Logger.WithFields(logrus.Fields{
			"rule":        "5",
			"item_index":  i + 1,
			"description": trimmedDesc,
			"length":      trimmedLen,
		})

		if trimmedLen > 0 && trimmedLen%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			calculation := price * 0.2
			itemPoints = int64(math.Ceil(calculation))
			totalPoints += itemPoints

			itemLogger.WithFields(logrus.Fields{
				"is_multiple_of_3": true,
				"price":            item.Price,
				"calculation":      calculation,
				"points":           itemPoints,
			}).Debug("Applied description length rule to item")
		} else {
			itemLogger.WithFields(logrus.Fields{
				"is_multiple_of_3": false,
				"points":           0,
			}).Debug("Applied description length rule to item")
		}
	}

	utils.Logger.WithFields(logrus.Fields{
		"rule":             "5",
		"points_from_rule": totalPoints,
	}).Debug("Applied description length rule total")

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

	utils.Logger.WithFields(logrus.Fields{
		"rule":             "7",
		"purchase_day":     day,
		"is_odd":           isOddDay,
		"points_from_rule": points,
	}).Debug("Applied odd day rule")

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

	utils.Logger.WithFields(logrus.Fields{
		"rule":                "8",
		"purchase_time":       purchaseTime,
		"is_in_special_range": inTimeRange,
		"points_from_rule":    points,
	}).Debug("Applied time range rule")

	return points
}
