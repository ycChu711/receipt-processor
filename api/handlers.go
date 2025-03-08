package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/ycChu711/receipt-processor/models"
	"github.com/ycChu711/receipt-processor/services"
	"github.com/ycChu711/receipt-processor/utils"
)

// ReceiptHandler handles receipt-related requests
type ReceiptHandler struct {
	service *services.ReceiptService
}

// NewReceiptHandler creates a new receipt handler
func NewReceiptHandler(service *services.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{
		service: service,
	}
}

// ProcessReceipt handles the POST /receipts/process endpoint
func (h *ReceiptHandler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	utils.Logger.Info("Processing receipt request")

	var receipt models.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		utils.Logger.WithError(err).Error("Failed to decode receipt JSON")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid receipt format. Please verify input."})
		return
	}

	// Validate the receipt
	if err := receipt.Validate(); err != nil {
		utils.Logger.WithError(err).Warn("Receipt validation failed")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid receipt. Please verify input. " + err.Error()})
		return
	}

	// Process the receipt
	id, err := h.service.ProcessReceipt(receipt)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to process receipt")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to process receipt"})
		return
	}

	// Return the ID
	utils.Logger.WithField("id", id).Info("Receipt processed successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.ReceiptResponse{ID: id})
}

// GetPoints handles the GET /receipts/{id}/points endpoint
func (h *ReceiptHandler) GetPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	utils.Logger.WithField("id", id).Info("Getting points for receipt")

	points, exists := h.service.GetPoints(id)
	if !exists {
		utils.Logger.WithField("id", id).Warn("Receipt not found")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No receipt found for that ID"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"id":     id,
		"points": points,
	}).Info("Points retrieved successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.PointsResponse{Points: points})
}
