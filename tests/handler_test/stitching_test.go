package handler_test

import (
	"Test-Task-Go/internal/handler"
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockStitchingService struct {
	ProcessStitchingFunc func(sourceID int, playlist string) (string, error)
}

func (m *mockStitchingService) ProcessStitching(ctx context.Context, sourceID int, playlist string) (string, error) {
	return m.ProcessStitchingFunc(sourceID, playlist)
}

func TestHandleStitching(t *testing.T) {
	tests := []struct {
		name                string
		sourceID            string
		mockPlaylist        []byte
		mockResponse        string
		mockError           error
		expectedStatusCode  int
		expectedContentType string
	}{
		{
			name:                "Invalid query parameters (sourceID=0)",
			sourceID:            "0",
			mockPlaylist:        []byte("#EXT-X-CUE-OUT:DURATION=30"),
			mockResponse:        `{"error": "Invalid query format. Ensure 'sourceID' is correctly provided."}`,
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: "application/json",
		},
		{
			name:                "No playlist provided",
			sourceID:            "1",
			mockPlaylist:        nil,
			mockResponse:        `{"error": "Playlist not provided"}`,
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: "application/json",
		},
		{
			name:                "Error in stitching service",
			sourceID:            "1",
			mockPlaylist:        []byte("#EXT-X-CUE-OUT:DURATION=30"),
			mockResponse:        `{"error": "Error processing playlist"}`,
			mockError:           errors.New("some error"),
			expectedStatusCode:  http.StatusInternalServerError,
			expectedContentType: "application/json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &mockStitchingService{
				ProcessStitchingFunc: func(sourceID int, playlist string) (string, error) {
					return tc.mockResponse, tc.mockError
				},
			}

			stitchingHandler := handler.NewStitchingHandler(mockService)
			req := httptest.NewRequest(http.MethodPost, "/stitching?sourceID="+tc.sourceID, bytes.NewReader(tc.mockPlaylist))
			w := httptest.NewRecorder()

			stitchingHandler.HandleStitching(w, req)
			res := w.Result()
			defer func(Body io.ReadCloser) {
				if err := Body.Close(); err != nil {
					log.Printf("Failed to close response body: %v", err)
				}
			}(res.Body)

			if res.StatusCode != tc.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tc.expectedStatusCode, res.StatusCode)
			}

			if res.Header.Get("Content-Type") != tc.expectedContentType {
				t.Errorf("expected content type %s, got %s", tc.expectedContentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
