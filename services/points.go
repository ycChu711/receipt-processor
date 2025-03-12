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
	var points int64 = 0

	utils.Logger.WithField("retailer", receipt.Retailer).Info("Starting points calculation")

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

	utils.Logger.WithFields(logrus.Fields{
		"rule":             "1",
		"retailer":         receipt.Retailer,
		"alphanumeric":     alphaNumChars,
		"count":            alphaNumCount,
		"points_from_rule": alphaNumCount,
		"running_total":    points,
	}).Debug("Applied retailer name points rule")

	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	totalFloat, _ := strconv.ParseFloat(receipt.Total, 64)
	rule2Points := int64(0)
	if totalFloat == float64(int64(totalFloat)) {
		rule2Points = 50
		points += rule2Points
		utils.Logger.WithFields(logrus.Fields{
			"rule":             "2",
			"total":            receipt.Total,
			"is_round_dollar":  true,
			"points_from_rule": rule2Points,
			"running_total":    points,
		}).Debug("Applied round dollar amount rule")
	} else {
		utils.Logger.WithFields(logrus.Fields{
			"rule":             "2",
			"total":            receipt.Total,
			"is_round_dollar":  false,
			"points_from_rule": rule2Points,
			"running_total":    points,
		}).Debug("Applied round dollar amount rule")
	}

	// Rule 3: 25 points if the total is a multiple of 0.25.
	rule3Points := int64(0)
	if math.Mod(totalFloat*100, 25) == 0 {
		rule3Points = 25
		points += rule3Points
		utils.Logger.WithFields(logrus.Fields{
			"rule":               "3",
			"total":              receipt.Total,
			"is_multiple_of_025": true,
			"points_from_rule":   rule3Points,
			"running_total":      points,
		}).Debug("Applied multiple of 0.25 rule")
	} else {
		utils.Logger.WithFields(logrus.Fields{
			"rule":               "3",
			"total":              receipt.Total,
			"is_multiple_of_025": false,
			"points_from_rule":   rule3Points,
			"running_total":      points,
		}).Debug("Applied multiple of 0.25 rule")
	}

	// Rule 4: 5 points for every two items on the receipt.
	itemPairs := len(receipt.Items) / 2
	pairPoints := int64(itemPairs * 5)
	points += pairPoints

	utils.Logger.WithFields(logrus.Fields{
		"rule":             "4",
		"items_count":      len(receipt.Items),
		"pairs":            itemPairs,
		"points_from_rule": pairPoints,
		"running_total":    points,
	}).Debug("Applied item pairs rule")

	// Rule 5: If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	var descPoints int64 = 0
	for i, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		trimmedLen := len(trimmedDesc)

		itemLogger := utils.Logger.WithFields(logrus.Fields{
			"rule":        "5",
			"item_index":  i + 1,
			"description": trimmedDesc,
			"length":      trimmedLen,
		})

		if trimmedLen > 0 && trimmedLen%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			calculation := price * 0.2
			itemPoints := int64(math.Ceil(calculation))
			descPoints += itemPoints

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
	points += descPoints

	utils.Logger.WithFields(logrus.Fields{
		"rule":             "5",
		"points_from_rule": descPoints,
		"running_total":    points,
	}).Debug("Applied description length rule total")

	// Rule 7: 6 points if the day in the purchase date is odd.
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	day := purchaseDate.Day()
	rule7Points := int64(0)

	if day%2 == 1 {
		rule7Points = 6
		points += rule7Points
		utils.Logger.WithFields(logrus.Fields{
			"rule":             "7",
			"purchase_day":     day,
			"is_odd":           true,
			"points_from_rule": rule7Points,
			"running_total":    points,
		}).Debug("Applied odd day rule")
	} else {
		utils.Logger.WithFields(logrus.Fields{
			"rule":             "7",
			"purchase_day":     day,
			"is_odd":           false,
			"points_from_rule": rule7Points,
			"running_total":    points,
		}).Debug("Applied odd day rule")
	}

	// Rule 8: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	hour := purchaseTime.Hour()
	minute := purchaseTime.Minute()

	inTimeRange := (hour == 14 && minute >= 0) || (hour == 15)
	rule8Points := int64(0)

	if inTimeRange {
		rule8Points = 10
		points += rule8Points
		utils.Logger.WithFields(logrus.Fields{
			"rule":                "8",
			"purchase_time":       receipt.PurchaseTime,
			"is_in_special_range": true,
			"points_from_rule":    rule8Points,
			"running_total":       points,
		}).Debug("Applied time range rule")
	} else {
		utils.Logger.WithFields(logrus.Fields{
			"rule":                "8",
			"purchase_time":       receipt.PurchaseTime,
			"is_in_special_range": false,
			"points_from_rule":    rule8Points,
			"running_total":       points,
		}).Debug("Applied time range rule")
	}

	utils.Logger.WithFields(logrus.Fields{
		"retailer":     receipt.Retailer,
		"final_points": points,
	}).Info("Completed points calculation")

	return points
}
