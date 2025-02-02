package handler_test

import (
	"Test-Task-Go/internal/entity"
	"Test-Task-Go/internal/handler"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuctionService struct {
	mock.Mock
}

func (m *MockAuctionService) ProcessAuction(ctx context.Context, sourceID, maxDuration int) ([]entity.Creative, error) {
	args := m.Called(ctx, sourceID, maxDuration)
	if creatives, ok := args.Get(0).([]entity.Creative); ok {
		return creatives, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestAuctionHandler_HandleAuction(t *testing.T) {
	tests := []struct {
		name               string
		sourceID           int
		maxDuration        int
		mockCreatives      []entity.Creative
		mockError          error
		expectedStatusCode int
		expectedBody       string
		expectProcessCall  bool
	}{
		{
			name:               "Invalid query parameters (sourceID=0)",
			sourceID:           0,
			maxDuration:        30,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"The sourceID value must be a positive integer."}`,
			expectProcessCall:  false,
		},
		{
			name:        "Successful auction",
			sourceID:    1,
			maxDuration: 30,
			mockCreatives: []entity.Creative{
				{ID: 1, CampaignID: 1, Duration: 30, Price: 5.00, PlaylistHLS: "#EXTINF:5.000, ad31.ts #EXTINF:5.000, " +
					"ad32.ts #EXTINF:5.000, ad33.ts #EXTINF:5.000, ad34.ts #EXTINF:5.000, ad35.ts"},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"creatives":[{"id":1,"campaign_id":1,"duration":30,"price":5.00,"playlist_hls":"#EXTINF:5.000, ad31.ts #EXTINF:5.000, ad32.ts #EXTINF:5.000, ad33.ts #EXTINF:5.000, ad34.ts #EXTINF:5.000, ad35.ts"}]}`,
			expectProcessCall:  true,
		},
		{
			name:               "Error processing auction",
			sourceID:           1,
			maxDuration:        30,
			mockCreatives:      nil,
			mockError:          errors.New("unable to process auction"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"unable to process auction"}`,
			expectProcessCall:  true,
		},
		{
			name:               "No creatives found",
			sourceID:           1,
			maxDuration:        30,
			mockCreatives:      nil,
			mockError:          errors.New("no creatives found for SourceID 1 with CueOutDuration 30"),
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       `{"error":"No creatives found for SourceID=1 with MaxDuration=30"}`,
			expectProcessCall:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuctionService := new(MockAuctionService)

			if tt.expectProcessCall {
				mockAuctionService.On("ProcessAuction", mock.Anything, tt.sourceID, tt.maxDuration).
					Return(tt.mockCreatives, tt.mockError)
			}

			auctionHandler := handler.NewAuctionHandler(mockAuctionService)

			req, _ := http.NewRequest("GET", "/auction?sourceID="+strconv.Itoa(tt.sourceID)+"&maxDuration="+strconv.Itoa(tt.maxDuration), nil)
			recorder := httptest.NewRecorder()

			auctionHandler.HandleAuction(recorder, req)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
			assert.JSONEq(t, tt.expectedBody, recorder.Body.String())

			if tt.expectProcessCall {
				mockAuctionService.AssertExpectations(t)
			} else {
				mockAuctionService.AssertNotCalled(t, "ProcessAuction", mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}
