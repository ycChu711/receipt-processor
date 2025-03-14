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

const (
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"
)

// ReceiptHandler manages HTTP requests for receipts
type ReceiptHandler struct {
	service *services.ReceiptService
}

// NewReceiptHandler creates a handler with the given service
func NewReceiptHandler(service *services.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{
		service: service,
	}
}

// ProcessReceipt handles POST /receipts/process
func (h *ReceiptHandler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	utils.Logger.Info("Processing receipt request")

	var receipt models.Receipt
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&receipt)

	// check bad json
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to decode receipt JSON")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set(headerContentType, contentTypeJSON)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid receipt format. Please verify input."})
		return
	}

	// validate receipt
	if err := receipt.Validate(); err != nil {
		utils.Logger.WithError(err).Warn("Receipt validation failed")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set(headerContentType, contentTypeJSON)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid receipt: " + err.Error(),
		})
		return
	}

	// process and get id
	id, err := h.service.ProcessReceipt(receipt)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to process receipt")
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set(headerContentType, contentTypeJSON)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Server error processing receipt",
		})
		return
	}

	// return id
	utils.Logger.WithField("id", id).Info("Receipt processed successfully")
	w.Header().Set(headerContentType, contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.ReceiptResponse{ID: id})
}

// GetPoints handles the GET /receipts/{id}/points
func (h *ReceiptHandler) GetPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	utils.Logger.WithField("id", id).Info("Getting points for receipt")

	points, found := h.service.GetPoints(id)
	if !found {
		utils.Logger.WithField("id", id).Warn("Receipt not found")
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set(headerContentType, contentTypeJSON)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "No receipt found for that ID",
		})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"id":     id,
		"points": points,
	}).Info("Got points for the receipt")

	w.Header().Set(headerContentType, contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.PointsResponse{Points: points})
}
