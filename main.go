package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ycChu711/receipt-processor/api"
	"github.com/ycChu711/receipt-processor/repository"
	"github.com/ycChu711/receipt-processor/services"
	"github.com/ycChu711/receipt-processor/utils"
)

func main() {
	// Intialize logger
	utils.InitLogger()
	utils.Logger.Info("Starting receipt processor API")

	// Create router
	r := mux.NewRouter()
	utils.Logger.Info("Router initialized")

	// Create storage
	storage := repository.NewInMemoryStorage()
	utils.Logger.Info("Storage initialized")

	// Create serive
	receiptService := services.NewReceiptService(storage)
	utils.Logger.Info("Receipt service initialized")

	// Setup routes
	api.SetupRoutes(r, receiptService)
	utils.Logger.Info("Routes configured")

	// Start server
	utils.Logger.Info("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		utils.Logger.WithError(err).Fatal("Server failed to start")
	}
}
