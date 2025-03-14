package models

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

func (r *Receipt) Validate() error {

	if err := validateRetailer(r.Retailer); err != nil {
		return err
	}

	if err := validateDateTime(r.PurchaseDate, r.PurchaseTime); err != nil {
		return err
	}

	if err := validateItems(r.Items); err != nil {
		return err
	}

	return validateTotal(r.Total)

}

func validateRetailer(retailer string) error {
	if strings.TrimSpace(retailer) == "" {
		return errors.New("Retailer is required")
	}
	retailerRegex := regexp.MustCompile(`^[\w\s\-&]+$`)
	if !retailerRegex.MatchString(retailer) {
		return errors.New("Retailer must be alphanumeric")
	}
	return nil
}

func validateDateTime(date, purchaseTime string) error {

	if date == "" {
		return errors.New("Purchase date is required")
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return errors.New("Invalid purchase date format, should be YYYY-MM-DD")
	}

	if purchaseTime == "" {
		return errors.New("Purchase time is required")
	}
	if _, err := time.Parse("15:04", purchaseTime); err != nil {
		return errors.New("Invalid purchase time format, should be HH:MM")
	}

	return nil
}

func validateItems(items []Item) error {
	if len(items) == 0 {
		return errors.New("Need at least one item")
	}

	for _, item := range items {
		if strings.TrimSpace(item.ShortDescription) == "" {
			return errors.New("Item description is required")
		}

		descRegex := regexp.MustCompile(`^[\w\s\-&]+$`)
		if !descRegex.MatchString(item.ShortDescription) {
			return errors.New("Item short description contains invalid characters")
		}

		if item.Price == "" {
			return errors.New("Item price is required")
		}

		priceRegex := regexp.MustCompile(`^\d+\.\d{2}$`)
		if !priceRegex.MatchString(item.Price) {
			return errors.New("Invalid item price format")
		}
	}
	return nil
}

func validateTotal(total string) error {
	if total == "" {
		return errors.New("Total is required")
	}

	totalRegex := regexp.MustCompile(`^\d+\.\d{2}$`)
	if !totalRegex.MatchString(total) {
		return errors.New("Invalid total format")
	}
	return nil
}
