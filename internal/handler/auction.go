package handler

import (
	"Test-Task-Go/internal/entity"
	"Test-Task-Go/internal/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type AuctionHandler struct {
	auctionService service.AuctionService
}

func NewAuctionHandler(auctionService service.AuctionService) *AuctionHandler {
	return &AuctionHandler{auctionService: auctionService}
}

func (h *AuctionHandler) HandleAuction(w http.ResponseWriter, r *http.Request) {
	sourceIDStr := r.URL.Query().Get("sourceID")
	maxDurationStr := r.URL.Query().Get("maxDuration")

	if sourceIDStr == "" || maxDurationStr == "" {
		log.Printf("Invalid query format: missing 'sourceID' or 'maxDuration'. SourceID=%s, MaxDuration=%s", sourceIDStr, maxDurationStr)
		w.Header().Set("Content-Type", "application/json")
		errorResponse := map[string]string{"error": "Invalid query format. Ensure 'sourceID' and 'maxDuration' are correctly provided."}
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			log.Printf("Failed to encode error response: %v", err)
		}
		return
	}

	sourceID, err := strconv.Atoi(sourceIDStr)
	if err != nil || sourceID <= 0 {
		log.Printf("Invalid 'sourceID' value: %s", sourceIDStr)
		w.Header().Set("Content-Type", "application/json")
		errorResponse := map[string]string{"error": "The sourceID value must be a positive integer."}
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			log.Printf("Failed to encode error response: %v", err)
		}
		return
	}

	maxDuration, err := strconv.Atoi(maxDurationStr)
	if err != nil || maxDuration <= 0 {
		log.Printf("Invalid 'maxDuration' value: %s", maxDurationStr)
		w.Header().Set("Content-Type", "application/json")
		errorResponse := map[string]string{"error": "The maxDuration value must be a positive integer."}
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			log.Printf("Failed to encode error response: %v", err)
		}
		return
	}

	log.Printf("Received auction request: SourceID=%d, MaxDuration=%d", sourceID, maxDuration)

	ctx := r.Context()

	creatives, err := h.auctionService.ProcessAuction(ctx, sourceID, maxDuration)
	if len(creatives) == 0 {
		log.Printf("Found 0 creatives for SourceID=%d with MaxDuration=%d", sourceID, maxDuration)
		w.Header().Set("Content-Type", "application/json")
		errorResponse := map[string]string{"error": fmt.Sprintf("No creatives found for SourceID=%d with MaxDuration=%d", sourceID, maxDuration)}
		w.WriteHeader(http.StatusNotFound)

		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			log.Printf("Failed to encode error response: %v", err)
		}
		return
	}

	if err != nil {
		if err.Error() == fmt.Sprintf("source with ID %d is inactive", sourceID) {
			log.Printf("Source with ID=%d is inactive", sourceID)
			w.Header().Set("Content-Type", "application/json")
			errorResponse := map[string]string{"error": fmt.Sprintf("Source with ID %d is inactive", sourceID)}
			w.WriteHeader(http.StatusNotFound)

			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				log.Printf("Failed to encode error response: %v", err)
			}
		} else {
			log.Printf("Error processing auction for SourceID=%d: %v", sourceID, err)
			w.Header().Set("Content-Type", "application/json")
			errorResponse := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusInternalServerError)

			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				log.Printf("Failed to encode error response: %v", err)
			}
		}
		return
	}

	log.Printf("Found %d creatives for SourceID=%d with MaxDuration=%d", len(creatives), sourceID, maxDuration)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := struct {
		Creatives []entity.Creative `json:"creatives"`
	}{
		Creatives: creatives,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}
