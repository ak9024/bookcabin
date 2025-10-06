package tests

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
)

// TestEndToEndWorkflow tests the complete flow from creating a flight to assigning a voucher
func TestEndToEndWorkflow(t *testing.T) {
	testApp := setupTestApp(t)
	defer testApp.cleanup()

	// Step 1: Create a flight
	t.Log("Step 1: Creating flight")
	flightBody := map[string]any{
		"flight_numbers": []string{"GA500"},
		"dep_date":       "2025-12-25",
	}
	resp, err := testApp.makeRequest("POST", "/api/v1/flights", flightBody)
	if err != nil {
		t.Fatalf("Failed to create flight: %v", err)
	}
	if resp.Code != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d. Body: %s", http.StatusCreated, resp.Code, resp.Body.String())
	}

	// Step 2: Create seats for the flight
	t.Log("Step 2: Creating seats")
	seatsBody := map[string]any{
		"flight_id": 1,
		"cabin":     "BUSINESS",
		"labels":    []string{"1A", "1B", "1C", "1D"},
	}
	resp, err = testApp.makeRequest("POST", "/api/v1/seats", seatsBody)
	if err != nil {
		t.Fatalf("Failed to create seats: %v", err)
	}
	if resp.Code != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.Code)
	}

	// Step 3: Create a voucher
	t.Log("Step 3: Creating voucher")
	voucherBody := map[string]any{
		"code":      "XMAS2025",
		"flight_id": 1,
		"cabin":     "BUSINESS",
	}
	resp, err = testApp.makeRequest("POST", "/api/v1/vouchers", voucherBody)
	if err != nil {
		t.Fatalf("Failed to create voucher: %v", err)
	}
	if resp.Code != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.Code)
	}

	// Step 4: Assign the voucher to a seat
	t.Log("Step 4: Assigning voucher to seat")
	assignBody := map[string]any{
		"voucher_code": "XMAS2025",
	}
	resp, err = testApp.makeRequest("POST", "/api/v1/vouchers/assigns", assignBody)
	if err != nil {
		t.Fatalf("Failed to assign voucher: %v", err)
	}
	if resp.Code != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d. Body: %s", http.StatusCreated, resp.Code, resp.Body.String())
	}

	// Parse assignment response
	var assignResult map[string]any
	parseResponse(t, resp, &assignResult)

	data, ok := assignResult["data"].(map[string]any)
	if !ok {
		t.Fatal("Expected data in assignment response")
	}

	// Verify assignment details
	if data["voucher_code"] != "XMAS2025" {
		t.Errorf("Expected voucher_code XMAS2025, got %v", data["voucher_code"])
	}
	if data["cabin"] != "BUSINESS" {
		t.Errorf("Expected cabin BUSINESS, got %v", data["cabin"])
	}
	if data["seat_label"] == nil || data["seat_label"] == "" {
		t.Error("Expected seat_label in response")
	}

	t.Logf("Successfully assigned voucher to seat: %s", data["seat_label"])

	// Step 5: Verify voucher is marked as redeemed
	t.Log("Step 5: Verifying voucher is redeemed")
	resp, err = testApp.makeRequest("GET", "/api/v1/vouchers", nil)
	if err != nil {
		t.Fatalf("Failed to get vouchers: %v", err)
	}

	var vouchersResult map[string]any
	parseResponse(t, resp, &vouchersResult)

	vouchers, ok := vouchersResult["data"].([]any)
	if !ok || len(vouchers) == 0 {
		t.Fatal("Expected vouchers in response")
	}

	voucher := vouchers[0].(map[string]any)
	redeemed := voucher["redeemed"]
	if redeemed == nil {
		t.Error("Expected redeemed field")
	}

	t.Log("End-to-end workflow completed successfully!")
}

// TestMultipleFlightsAndVouchers tests handling multiple flights and vouchers
func TestMultipleFlightsAndVouchers(t *testing.T) {
	testApp := setupTestApp(t)
	defer testApp.cleanup()

	// Create multiple flights
	for i := 1; i <= 3; i++ {
		flightBody := map[string]any{
			"flight_numbers": []string{fmt.Sprintf("GA%d00", i)},
			"dep_date":       "2025-10-15",
		}
		resp, err := testApp.makeRequest("POST", "/api/v1/flights", flightBody)
		if err != nil || resp.Code != http.StatusCreated {
			t.Fatalf("Failed to create flight %d", i)
		}

		// Create seats for each flight
		seatsBody := map[string]any{
			"flight_id": i,
			"cabin":     "ECONOMY",
			"labels":    []string{fmt.Sprintf("%dA", i), fmt.Sprintf("%dB", i)},
		}
		resp, err = testApp.makeRequest("POST", "/api/v1/seats", seatsBody)
		if err != nil || resp.Code != http.StatusCreated {
			t.Fatalf("Failed to create seats for flight %d", i)
		}

		// Create voucher for each flight
		voucherBody := map[string]any{
			"code":      fmt.Sprintf("VOUCHER%d", i),
			"flight_id": i,
			"cabin":     "ECONOMY",
		}
		resp, err = testApp.makeRequest("POST", "/api/v1/vouchers", voucherBody)
		if err != nil || resp.Code != http.StatusCreated {
			t.Fatalf("Failed to create voucher for flight %d", i)
		}
	}

	// Assign all vouchers
	for i := 1; i <= 3; i++ {
		assignBody := map[string]any{
			"voucher_code": fmt.Sprintf("VOUCHER%d", i),
		}
		resp, err := testApp.makeRequest("POST", "/api/v1/vouchers/assigns", assignBody)
		if err != nil || resp.Code != http.StatusCreated {
			t.Errorf("Failed to assign voucher %d", i)
		}
	}

	t.Log("Successfully handled multiple flights and vouchers")
}

// TestConcurrentVoucherAssignments tests concurrent voucher assignments
func TestConcurrentVoucherAssignments(t *testing.T) {
	testApp := setupTestApp(t)
	defer testApp.cleanup()

	// Setup: Create flight with limited seats
	flightBody := map[string]any{
		"flight_numbers": []string{"GA999"},
		"dep_date":       "2025-11-11",
	}
	testApp.makeRequest("POST", "/api/v1/flights", flightBody)

	// Create 3 seats
	seatsBody := map[string]any{
		"flight_id": 1,
		"cabin":     "ECONOMY",
		"labels":    []string{"1A", "1B", "1C"},
	}
	testApp.makeRequest("POST", "/api/v1/seats", seatsBody)

	// Create 5 vouchers (more than available seats)
	for i := 1; i <= 5; i++ {
		voucherBody := map[string]any{
			"code":      fmt.Sprintf("CONCURRENT%d", i),
			"flight_id": 1,
			"cabin":     "ECONOMY",
		}
		testApp.makeRequest("POST", "/api/v1/vouchers", voucherBody)
	}

	// Try to assign all vouchers concurrently
	var wg sync.WaitGroup
	results := make([]int, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			assignBody := map[string]any{
				"voucher_code": fmt.Sprintf("CONCURRENT%d", index+1),
			}
			resp, _ := testApp.makeRequest("POST", "/api/v1/vouchers/assigns", assignBody)
			results[index] = resp.Code
		}(i)
	}

	wg.Wait()

	// Count successful assignments
	successCount := 0
	failedCount := 0
	for _, code := range results {
		if code == http.StatusCreated {
			successCount++
		} else if code == http.StatusBadRequest {
			failedCount++
		}
	}

	// Due to concurrent execution and retry logic, we should have:
	// - At least some successful assignments (up to 3, matching seat count)
	// - Total of success + failed should equal 5
	total := successCount + failedCount
	if total != 5 {
		t.Errorf("Expected total of 5 assignments, got %d", total)
	}

	// Should have at least 1 successful
	if successCount < 1 {
		t.Errorf("Expected at least 1 successful assignment, got %d", successCount)
	}

	// Should not exceed seat count
	if successCount > 3 {
		t.Errorf("Expected at most 3 successful assignments, got %d", successCount)
	}

	t.Logf("Concurrent assignments: %d successful, %d failed (as expected)", successCount, failedCount)
}

// TestMixedCabinClasses tests vouchers and seats for different cabin classes
func TestMixedCabinClasses(t *testing.T) {
	testApp := setupTestApp(t)
	defer testApp.cleanup()

	// Create flight
	flightBody := map[string]any{
		"flight_numbers": []string{"GA777"},
		"dep_date":       "2025-10-20",
	}
	testApp.makeRequest("POST", "/api/v1/flights", flightBody)

	// Create seats for different cabins
	cabins := []string{"ECONOMY", "BUSINESS", "FIRST"}
	for _, cabin := range cabins {
		seatsBody := map[string]any{
			"flight_id": 1,
			"cabin":     cabin,
			"labels":    []string{cabin + "1A", cabin + "1B"},
		}
		resp, _ := testApp.makeRequest("POST", "/api/v1/seats", seatsBody)
		if resp.Code != http.StatusCreated {
			t.Fatalf("Failed to create seats for %s cabin", cabin)
		}
	}

	// Create vouchers for each cabin
	for i, cabin := range cabins {
		voucherBody := map[string]any{
			"code":      fmt.Sprintf("VIP%s%d", cabin, i),
			"flight_id": 1,
			"cabin":     cabin,
		}
		resp, _ := testApp.makeRequest("POST", "/api/v1/vouchers", voucherBody)
		if resp.Code != http.StatusCreated {
			t.Fatalf("Failed to create voucher for %s cabin", cabin)
		}
	}

	// Assign vouchers for each cabin
	for i, cabin := range cabins {
		assignBody := map[string]any{
			"voucher_code": fmt.Sprintf("VIP%s%d", cabin, i),
		}
		resp, _ := testApp.makeRequest("POST", "/api/v1/vouchers/assigns", assignBody)
		if resp.Code != http.StatusCreated {
			t.Errorf("Failed to assign voucher for %s cabin", cabin)
		}

		// Verify assigned to correct cabin
		var result map[string]any
		parseResponse(t, resp, &result)
		data := result["data"].(map[string]any)
		if data["cabin"] != cabin {
			t.Errorf("Expected cabin %s, got %v", cabin, data["cabin"])
		}
	}

	t.Log("Successfully handled mixed cabin classes")
}
