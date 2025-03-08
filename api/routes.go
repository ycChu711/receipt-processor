package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ycChu711/receipt-processor/services"
)

// SetupRoutes configures the API routes
func SetupRoutes(r *mux.Router, receiptService *services.ReceiptService) {
	// Create a new receipt handler
	receiptHandler := NewReceiptHandler(receiptService)

	// Set up routes
	r.HandleFunc("/receipts/process", receiptHandler.ProcessReceipt).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", receiptHandler.GetPoints).Methods("GET")

	// Add health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

}
