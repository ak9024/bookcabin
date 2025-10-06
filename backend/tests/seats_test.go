package tests

import (
	"net/http"
	"testing"
)

func TestCreateSeats(t *testing.T) {
	testApp := setupTestApp(t)
	defer testApp.cleanup()

	flightBody := map[string]any{
		"flight_numbers": []string{"GA100"},
		"dep_date":       "2025-10-10",
	}
	_, err := testApp.makeRequest("POST", "/api/v1/flights", flightBody)
	if err != nil {
		t.Fatalf("Failed to create flight: %v", err)
	}

	tests := []struct {
		name           string
		requestBody    map[string]any
		expectedStatus int
	}{
		{
			name: "Create economy seats successfully",
			requestBody: map[string]any{
				"flight_id": 1,
				"cabin":     "ECONOMY",
				"labels":    []string{"1A", "1B", "1C"},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Create business seats successfully",
			requestBody: map[string]any{
				"flight_id": 1,
				"cabin":     "BUSINESS",
				"labels":    []string{"2A", "2B"},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid flight ID",
			requestBody: map[string]any{
				"flight_id": 999,
				"cabin":     "ECONOMY",
				"labels":    []string{"3A"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty labels",
			requestBody: map[string]any{
				"flight_id": 1,
				"cabin":     "ECONOMY",
				"labels":    []string{},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing flight_id",
			requestBody: map[string]any{
				"cabin":  "ECONOMY",
				"labels": []string{"4A"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Zero flight_id",
			requestBody: map[string]any{
				"flight_id": 0,
				"cabin":     "ECONOMY",
				"labels":    []string{"5A"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Negative flight_id",
			requestBody: map[string]any{
				"flight_id": -1,
				"cabin":     "ECONOMY",
				"labels":    []string{"6A"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing cabin field",
			requestBody: map[string]any{
				"flight_id": 1,
				"labels":    []string{"7A"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid cabin value",
			requestBody: map[string]any{
				"flight_id": 1,
				"cabin":     "PREMIUM",
				"labels":    []string{"8A"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty string in labels array",
			requestBody: map[string]any{
				"flight_id": 1,
				"cabin":     "ECONOMY",
				"labels":    []string{"9A", ""},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing labels field",
			requestBody: map[string]any{
				"flight_id": 1,
				"cabin":     "ECONOMY",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := testApp.makeRequest("POST", "/api/v1/seats", tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}

			if resp.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, resp.Code, resp.Body.String())
			}
		})
	}
}

func TestGetAllSeats(t *testing.T) {
	testApp := setupTestApp(t)
	defer testApp.cleanup()

	flightBody := map[string]any{
		"flight_numbers": []string{"GA100"},
		"dep_date":       "2025-10-10",
	}
	_, err := testApp.makeRequest("POST", "/api/v1/flights", flightBody)
	if err != nil {
		t.Fatalf("Failed to create flight: %v", err)
	}

	seatsBody := map[string]any{
		"flight_id": 1,
		"cabin":     "ECONOMY",
		"labels":    []string{"1A", "1B", "1C"},
	}
	_, err = testApp.makeRequest("POST", "/api/v1/seats", seatsBody)
	if err != nil {
		t.Fatalf("Failed to create seats: %v", err)
	}

	resp, err := testApp.makeRequest("GET", "/api/v1/seats", nil)
	if err != nil {
		t.Fatalf("Failed to get seats: %v", err)
	}

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
	}

	var result map[string]any
	parseResponse(t, resp, &result)

	if result["data"] == nil {
		t.Error("Expected data in response")
	}
}

func TestGetAllSeatsEmpty(t *testing.T) {
	testApp := setupTestApp(t)
	defer testApp.cleanup()

	resp, err := testApp.makeRequest("GET", "/api/v1/seats", nil)
	if err != nil {
		t.Fatalf("Failed to get seats: %v", err)
	}

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
	}
}
