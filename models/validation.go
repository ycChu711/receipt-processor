package models

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// Validate validates a receipt
func (r *Receipt) Validate() error {
	// Validate retailer
	if strings.TrimSpace(r.Retailer) == "" {
		return errors.New("Retailer is required")
	}
	retailerRegax := regexp.MustCompile(`^[\w\s\-&]+$`) // Alphanumeric, space, hyphen, and ampersand
	if !retailerRegax.MatchString(r.Retailer) {
		return errors.New("Retailer must be alphanumeric")
	}

	// Validate purchase date
	if r.PurchaseDate == "" {
		return errors.New("Purchase date is required")
	}
	_, err := time.Parse("2006-01-02", r.PurchaseDate) // YYYY-MM-DD
	if err != nil {
		return errors.New("Invalid purchase date format (should be YYYY-MM-DD)")
	}

	// Validate purchase time. 24-hour time expected
	if r.PurchaseTime == "" {
		return errors.New("Purchase time is required")
	}
	_, err = time.Parse("15:04", r.PurchaseTime) // HH:MM
	if err != nil {
		return errors.New("Invalid purchase time format (should be HH:MM)")
	}

	// Validate items
	if len(r.Items) == 0 {
		return errors.New("At least one item is required")
	}
	for _, item := range r.Items {
		if strings.TrimSpace(item.ShortDescription) == "" {
			return errors.New("Item short description is required")
		}

		descRegax := regexp.MustCompile(`^[\w\s\-&]+$`) // Alphanumeric, space, hyphen, and ampersand
		if !descRegax.MatchString(item.ShortDescription) {
			return errors.New("Item short description contains invalid characters")
		}

		if item.Price == "" {
			return errors.New("Item price is required")
		}

		priceRegax := regexp.MustCompile(`^\d+\.\d{2}$`) // 0.00
		if !priceRegax.MatchString(item.Price) {
			return errors.New("Invalid item price format")
		}

	}

	// Validate total
	if r.Total == "" {
		return errors.New("Total is required")
	}

	totalRegax := regexp.MustCompile(`^\d+\.\d{2}$`) // 0.00
	if !totalRegax.MatchString(r.Total) {
		return errors.New("Invalid total format")
	}

	return nil

}
