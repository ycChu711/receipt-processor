package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ycChu711/receipt-processor/services"
)

// SetupRoutes registers all API endpoints
func SetupRoutes(r *mux.Router, receiptService *services.ReceiptService) {
	receiptHandler := NewReceiptHandler(receiptService)

	r.HandleFunc("/receipts/process", receiptHandler.ProcessReceipt).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", receiptHandler.GetPoints).Methods("GET")

	// healthCheck responds with a simple status for monitoring
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

}
