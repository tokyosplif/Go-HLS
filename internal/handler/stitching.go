package handler

import (
	"Test-Task-Go/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type StitchingHandler struct {
	stitchingService service.StitchingService
}

func NewStitchingHandler(stitchingService service.StitchingService) *StitchingHandler {
	return &StitchingHandler{stitchingService: stitchingService}
}

func (h *StitchingHandler) HandleStitching(w http.ResponseWriter, r *http.Request) {
	sourceIDStr := r.URL.Query().Get("sourceID")
	sourceID, err := strconv.Atoi(sourceIDStr)
	if err != nil || sourceID <= 0 {
		log.Printf("Error parsing request: invalid sourceID")
		w.Header().Set("Content-Type", "application/json")
		errorResponse := map[string]string{"error": "Invalid query format. Ensure 'sourceID' is correctly provided."}
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			log.Printf("Failed to encode error response: %v", err)
		}
		return
	}

	rawPlaylist, err := io.ReadAll(r.Body)
	if err != nil || len(rawPlaylist) == 0 {
		w.Header().Set("Content-Type", "application/json")
		errorResponse := map[string]string{"error": "Playlist not provided"}
		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			log.Printf("Failed to encode error response: %v", err)
		}
		return
	}

	ctx := context.Background()

	modifiedPlaylist, err := h.stitchingService.ProcessStitching(ctx, sourceID, string(rawPlaylist))
	if err != nil {
		if err.Error() == fmt.Sprintf("no creatives found for SourceID %d with CueOutDuration %d", sourceID, 0) { // Adjust maxDuration if needed
			log.Printf("No creatives found for SourceID=%d", sourceID)
			w.Header().Set("Content-Type", "application/json")
			errorResponse := map[string]string{"error": fmt.Sprintf("No creatives found for SourceID=%d", sourceID)}
			w.WriteHeader(http.StatusNotFound)

			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				log.Printf("Failed to encode error response: %v", err)
			}
			return
		}

		if err.Error() == fmt.Sprintf("source with ID %d is inactive", sourceID) {
			w.Header().Set("Content-Type", "application/json")
			errorResponse := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusNotFound)

			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				log.Printf("Failed to encode error response: %v", err)
			}
			return
		}

		log.Printf("Error processing playlist for SourceID=%d: %v", sourceID, err)
		w.Header().Set("Content-Type", "application/json")
		errorResponse := map[string]string{"error": "Error processing playlist"}
		w.WriteHeader(http.StatusInternalServerError)

		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			log.Printf("Failed to encode error response: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(modifiedPlaylist))
}
