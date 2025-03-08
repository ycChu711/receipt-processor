package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ycChu711/receipt-processor/models"
	"github.com/ycChu711/receipt-processor/repository"
	"github.com/ycChu711/receipt-processor/services"
)

func setupTestHandler() *ReceiptHandler {
	storage := repository.NewInMemoryStorage()
	service := services.NewReceiptService(storage)
	return NewReceiptHandler(service)
}

func TestProcessReceipt(t *testing.T) {
	// Test valid receipt
	handler := setupTestHandler()

	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
		},
		Total: "6.49",
	}

	body, _ := json.Marshal(receipt)
	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	http.HandlerFunc(handler.ProcessReceipt).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.ReceiptResponse
	json.Unmarshal(rr.Body.Bytes(), &response)
	if response.ID == "" {
		t.Errorf("Expected non-empty ID")
	}

	// Test invalid receipt (missing retailer)
	invalidReceipt := models.Receipt{
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Item", Price: "1.00"},
		},
		Total: "1.00",
	}

	body, _ = json.Marshal(invalidReceipt)
	req, _ = http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	http.HandlerFunc(handler.ProcessReceipt).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetPoints(t *testing.T) {
	// Setup test data
	storage := repository.NewInMemoryStorage()
	service := services.NewReceiptService(storage)
	handler := NewReceiptHandler(service)

	// Process a receipt first to get a valid ID
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Item", Price: "1.00"},
		},
		Total: "1.00",
	}

	body, _ := json.Marshal(receipt)
	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	http.HandlerFunc(handler.ProcessReceipt).ServeHTTP(rr, req)

	var processResponse models.ReceiptResponse
	json.Unmarshal(rr.Body.Bytes(), &processResponse)
	validID := processResponse.ID

	// Test with valid ID
	req, _ = http.NewRequest("GET", "/receipts/"+validID+"/points", nil)
	vars := map[string]string{
		"id": validID,
	}
	req = mux.SetURLVars(req, vars)

	rr = httptest.NewRecorder()
	http.HandlerFunc(handler.GetPoints).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var pointsResponse models.PointsResponse
	json.Unmarshal(rr.Body.Bytes(), &pointsResponse)
	if pointsResponse.Points != 60 { // 6 points for Target + 50 points for round dollar + 4 points for odd day
		t.Errorf("Expected 60 points, got %v", pointsResponse.Points)
	}

	// Test with invalid ID
	req, _ = http.NewRequest("GET", "/receipts/invalid-id/points", nil)
	vars = map[string]string{
		"id": "invalid-id",
	}
	req = mux.SetURLVars(req, vars)

	rr = httptest.NewRecorder()
	http.HandlerFunc(handler.GetPoints).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
