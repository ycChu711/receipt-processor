package models

// Receipt represents a receipt of a transaction
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

// Item represents an item in the receipt
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// ReceiptResponse is the response for processing a receipt
type ReceiptResponse struct {
	ID string `json:"id"`
}

// PointsResponse is the response for getting points
type PointsResponse struct {
	Points int64 `json:"points"`
}
