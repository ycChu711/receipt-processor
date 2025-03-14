package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ycChu711/receipt-processor/models"
	"github.com/ycChu711/receipt-processor/repository"
	"github.com/ycChu711/receipt-processor/services"
)

const (
	testDate          = "2022-01-01"
	testTime          = "13:01"
	processEndpoint   = "/receipts/process"
	contentTypeHeader = "Content-Type"
	jsonContentType   = "application/json"
)

func createTestHandler() *ReceiptHandler {
	return NewReceiptHandler(
		services.NewReceiptService(
			repository.NewInMemoryStorage(),
		),
	)
}

func TestProcessReceipt(t *testing.T) {

	handler := createTestHandler()

	// valid receipt
	t.Run("valid receipt", func(t *testing.T) {
		// Make a simple test receipt
		receipt := models.Receipt{
			Retailer:     "Safeway",
			PurchaseDate: testDate,
			PurchaseTime: testTime,
			Items: []models.Item{
				{ShortDescription: "Ice Cream", Price: "5.99"},
			},
			Total: "5.99",
		}

		response := sendPostRequest(t, handler.ProcessReceipt, processEndpoint, receipt)

		if response.Code != http.StatusOK {
			t.Fatalf("Should get 200 OK but got %d instead", response.Code)
		}

		var respData models.ReceiptResponse
		json.Unmarshal(response.Body.Bytes(), &respData)

		if respData.ID == "" {
			t.Fatal("not getting a valid id")
		}
	})

	// Missing retailer
	t.Run("missing retailer", func(t *testing.T) {
		badReceipt := models.Receipt{
			// Missing retailer
			PurchaseDate: testDate,
			PurchaseTime: testTime,
			Items: []models.Item{
				{ShortDescription: "Candy", Price: "1.25"},
			},
			Total: "1.25",
		}

		response := sendPostRequest(t, handler.ProcessReceipt, processEndpoint, badReceipt)

		if response.Code != http.StatusBadRequest {
			t.Fatalf("Should get 400 for invalid receipt, got %d", response.Code)
		}
	})
}

func TestGetPoints(t *testing.T) {
	storage := repository.NewInMemoryStorage()
	service := services.NewReceiptService(storage)
	handler := NewReceiptHandler(service)

	receipt := models.Receipt{
		Retailer:     "Shop",
		PurchaseDate: testDate,
		PurchaseTime: testTime,
		Items: []models.Item{
			{ShortDescription: "Item", Price: "2.49"},
			{ShortDescription: "Coke", Price: "3.29"},
		},
		Total: "5.78",
	}

	response := sendPostRequest(t, handler.ProcessReceipt, processEndpoint, receipt)
	var processResp models.ReceiptResponse
	json.Unmarshal(response.Body.Bytes(), &processResp)
	validID := processResp.ID

	// valid id
	t.Run("valid receipt ID", func(t *testing.T) {

		req, _ := http.NewRequest("GET", fmt.Sprintf("/receipts/%s/points", validID), nil)
		req = mux.SetURLVars(req, map[string]string{"id": validID})

		recorder := httptest.NewRecorder()
		handler.GetPoints(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("Should get 200 OK but got %d", recorder.Code)
		}

		var pointsResp models.PointsResponse
		json.Unmarshal(recorder.Body.Bytes(), &pointsResp)

		if pointsResp.Points <= 0 {
			t.Errorf("Invalid points, got %d", pointsResp.Points)
		}
	})

	// non-exist id
	t.Run("nonexistent receipt ID", func(t *testing.T) {
		fakeID := "non-exist-id"
		req, _ := http.NewRequest("GET", fmt.Sprintf("/receipts/%s/points", fakeID), nil)
		req = mux.SetURLVars(req, map[string]string{"id": fakeID})

		recorder := httptest.NewRecorder()
		handler.GetPoints(recorder, req)

		if recorder.Code != http.StatusNotFound {
			t.Fatalf("Should get 404 Not Found for nonexistent ID, got %d", recorder.Code)
		}
	})
}

// Helper function to send POST requests and return the response
func sendPostRequest(t *testing.T, handlerFunc http.HandlerFunc, endpoint string, data interface{}) *httptest.ResponseRecorder {
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set(contentTypeHeader, jsonContentType)

	recorder := httptest.NewRecorder()
	handlerFunc(recorder, req)

	return recorder
}
