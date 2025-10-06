package tests

import (
	"net/http"
	"testing"
)

func TestCreateFlights(t *testing.T) {
	testApp := setupTestApp(t)
	defer testApp.cleanup()

	tests := []struct {
		name           string
		requestBody    map[string]any
		expectedStatus int
		checkResponse  func(t *testing.T, statusCode int, body string)
	}{
		{
			name: "Create flights successfully",
			requestBody: map[string]any{
				"flight_numbers": []string{"GA100", "GA200"},
				"dep_date":       "2025-10-10",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, statusCode int, body string) {
				if statusCode != http.StatusCreated {
					t.Errorf("Expected status %d, got %d", http.StatusCreated, statusCode)
				}
			},
		},
		{
			name: "Invalid date format",
			requestBody: map[string]any{
				"flight_numbers": []string{"GA300"},
				"dep_date":       "10-10-2025", // wrong format
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, statusCode int, body string) {
				if statusCode != http.StatusBadRequest {
					t.Errorf("Expected status %d, got %d", http.StatusBadRequest, statusCode)
				}
			},
		},
		{
			name: "Empty flight numbers",
			requestBody: map[string]any{
				"flight_numbers": []string{},
				"dep_date":       "2025-10-10",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, statusCode int, body string) {
				if statusCode != http.StatusBadRequest {
					t.Errorf("Expected status %d, got %d", http.StatusBadRequest, statusCode)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := testApp.makeRequest("POST", "/api/v1/flights", tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}

			if resp.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, resp.Code, resp.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, resp.Code, resp.Body.String())
			}
		})
	}
}

func TestGetAllFlights(t *testing.T) {
	testApp := setupTestApp(t)
	defer testApp.cleanup()

	// First, create some flights
	createBody := map[string]any{
		"flight_numbers": []string{"GA100", "GA200"},
		"dep_date":       "2025-10-10",
	}
	_, err := testApp.makeRequest("POST", "/api/v1/flights", createBody)
	if err != nil {
		t.Fatalf("Failed to create flights: %v", err)
	}

	// Now get all flights
	resp, err := testApp.makeRequest("GET", "/api/v1/flights", nil)
	if err != nil {
		t.Fatalf("Failed to get flights: %v", err)
	}

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
	}

	// Parse response
	var result map[string]any
	parseResponse(t, resp, &result)

	// Check that we have data
	if result["data"] == nil {
		t.Error("Expected data in response")
	}
}

func TestGetAllFlightsEmpty(t *testing.T) {
	testApp := setupTestApp(t)
	defer testApp.cleanup()

	resp, err := testApp.makeRequest("GET", "/api/v1/flights", nil)
	if err != nil {
		t.Fatalf("Failed to get flights: %v", err)
	}

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
	}
}
