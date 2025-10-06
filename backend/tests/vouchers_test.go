package tests

import (
	"net/http"
	"testing"
)

func TestCreateVoucher(t *testing.T) {
	testApp := setupTestApp(t)
	defer testApp.cleanup()

	// Setup: Create flight and seats
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
		"labels":    []string{"1A", "1B"},
	}
	_, err = testApp.makeRequest("POST", "/api/v1/seats", seatsBody)
	if err != nil {
		t.Fatalf("Failed to create seats: %v", err)
	}

	tests := []struct {
		name           string
		requestBody    map[string]any
		expectedStatus int
	}{
		{
			name: "Create voucher successfully",
			requestBody: map[string]any{
				"code":      "VOUCHER100",
				"flight_id": 1,
				"cabin":     "ECONOMY",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid flight ID",
			requestBody: map[string]any{
				"code":      "VOUCHER200",
				"flight_id": 999,
				"cabin":     "ECONOMY",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "No seats available for cabin",
			requestBody: map[string]any{
				"code":      "VOUCHER300",
				"flight_id": 1,
				"cabin":     "BUSINESS", // No business seats created
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := testApp.makeRequest("POST", "/api/v1/vouchers", tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}

			if resp.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, resp.Code, resp.Body.String())
			}
		})
	}
}

func TestGetAllVouchers(t *testing.T) {
	testApp := setupTestApp(t)
	defer testApp.cleanup()

	// Setup: Create flight, seats, and voucher
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
		"labels":    []string{"1A"},
	}
	_, err = testApp.makeRequest("POST", "/api/v1/seats", seatsBody)
	if err != nil {
		t.Fatalf("Failed to create seats: %v", err)
	}

	voucherBody := map[string]any{
		"code":      "VOUCHER100",
		"flight_id": 1,
		"cabin":     "ECONOMY",
	}
	_, err = testApp.makeRequest("POST", "/api/v1/vouchers", voucherBody)
	if err != nil {
		t.Fatalf("Failed to create voucher: %v", err)
	}

	// Get all vouchers
	resp, err := testApp.makeRequest("GET", "/api/v1/vouchers", nil)
	if err != nil {
		t.Fatalf("Failed to get vouchers: %v", err)
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

func TestAssignVoucher(t *testing.T) {
	testApp := setupTestApp(t)
	defer testApp.cleanup()

	// Setup: Create flight, seats, and voucher
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

	voucherBody := map[string]any{
		"code":      "VOUCHER100",
		"flight_id": 1,
		"cabin":     "ECONOMY",
	}
	_, err = testApp.makeRequest("POST", "/api/v1/vouchers", voucherBody)
	if err != nil {
		t.Fatalf("Failed to create voucher: %v", err)
	}

	tests := []struct {
		name           string
		requestBody    map[string]any
		expectedStatus int
	}{
		{
			name: "Assign voucher successfully",
			requestBody: map[string]any{
				"voucher_code": "VOUCHER100",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Voucher not found",
			requestBody: map[string]any{
				"voucher_code": "INVALID",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := testApp.makeRequest("POST", "/api/v1/vouchers/assigns", tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}

			if resp.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, resp.Code, resp.Body.String())
			}

			// For successful assignment, verify response contains seat assignment
			if tt.expectedStatus == http.StatusCreated {
				var result map[string]any
				parseResponse(t, resp, &result)

				data, ok := result["data"].(map[string]any)
				if !ok {
					t.Error("Expected data object in response")
					return
				}

				if data["seat_label"] == nil {
					t.Error("Expected seat_label in response")
				}
				if data["cabin"] == nil {
					t.Error("Expected cabin in response")
				}
			}
		})
	}
}

func TestAssignVoucherAlreadyRedeemed(t *testing.T) {
	testApp := setupTestApp(t)
	defer testApp.cleanup()

	// Setup: Create flight, seats, and voucher
	flightBody := map[string]any{
		"flight_numbers": []string{"GA100"},
		"dep_date":       "2025-10-10",
	}
	testApp.makeRequest("POST", "/api/v1/flights", flightBody)

	seatsBody := map[string]any{
		"flight_id": 1,
		"cabin":     "ECONOMY",
		"labels":    []string{"1A", "1B"},
	}
	testApp.makeRequest("POST", "/api/v1/seats", seatsBody)

	voucherBody := map[string]any{
		"code":      "VOUCHER100",
		"flight_id": 1,
		"cabin":     "ECONOMY",
	}
	testApp.makeRequest("POST", "/api/v1/vouchers", voucherBody)

	// Assign voucher first time
	assignBody := map[string]any{
		"voucher_code": "VOUCHER100",
	}
	resp, _ := testApp.makeRequest("POST", "/api/v1/vouchers/assigns", assignBody)
	if resp.Code != http.StatusCreated {
		t.Fatalf("First assignment should succeed, got status %d", resp.Code)
	}

	// Try to assign again - should fail
	resp, err := testApp.makeRequest("POST", "/api/v1/vouchers/assigns", assignBody)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for already redeemed voucher, got %d", http.StatusBadRequest, resp.Code)
	}
}

func TestAssignVoucherNoSeatsAvailable(t *testing.T) {
	testApp := setupTestApp(t)
	defer testApp.cleanup()

	// Setup: Create flight with only 1 seat
	flightBody := map[string]any{
		"flight_numbers": []string{"GA100"},
		"dep_date":       "2025-10-10",
	}
	testApp.makeRequest("POST", "/api/v1/flights", flightBody)

	seatsBody := map[string]any{
		"flight_id": 1,
		"cabin":     "ECONOMY",
		"labels":    []string{"1A"}, // Only 1 seat
	}
	testApp.makeRequest("POST", "/api/v1/seats", seatsBody)

	// Create 2 vouchers
	voucher1Body := map[string]any{
		"code":      "VOUCHER1",
		"flight_id": 1,
		"cabin":     "ECONOMY",
	}
	testApp.makeRequest("POST", "/api/v1/vouchers", voucher1Body)

	voucher2Body := map[string]any{
		"code":      "VOUCHER2",
		"flight_id": 1,
		"cabin":     "ECONOMY",
	}
	testApp.makeRequest("POST", "/api/v1/vouchers", voucher2Body)

	// Assign first voucher - should succeed
	assign1 := map[string]any{
		"voucher_code": "VOUCHER1",
	}
	resp1, _ := testApp.makeRequest("POST", "/api/v1/vouchers/assigns", assign1)
	if resp1.Code != http.StatusCreated {
		t.Fatalf("First assignment should succeed, got status %d", resp1.Code)
	}

	// Assign second voucher - should fail (no seats available)
	assign2 := map[string]any{
		"voucher_code": "VOUCHER2",
	}
	resp2, _ := testApp.makeRequest("POST", "/api/v1/vouchers/assigns", assign2)
	if resp2.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d when no seats available, got %d", http.StatusBadRequest, resp2.Code)
	}
}
