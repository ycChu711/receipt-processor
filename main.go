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
	utils.InitLogger()
	utils.Logger.Info("Starting receipt processor service...")

	// create router
	r := mux.NewRouter()

	// create storage and service
	storage := repository.NewInMemoryStorage()
	receiptService := services.NewReceiptService(storage)

	// setup api routes
	api.SetupRoutes(r, receiptService)
	utils.Logger.Info("API Routes configured")

	// Start server
	utils.Logger.Info("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		utils.Logger.WithError(err).Fatal("Server failed to start")
	}
}
