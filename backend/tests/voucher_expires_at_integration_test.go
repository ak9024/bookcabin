package tests

import (
	"net/http"
	"testing"
)

func TestVoucherExpiresAtDatabase(t *testing.T) {
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
		"labels":    []string{"1A", "1B"},
	}
	_, err = testApp.makeRequest("POST", "/api/v1/seats", seatsBody)
	if err != nil {
		t.Fatalf("Failed to create seats: %v", err)
	}

	t.Run("Voucher without expires_at should have NULL in database", func(t *testing.T) {
		voucherBody := map[string]any{
			"code":      "VOUCHER_NO_EXPIRY",
			"flight_id": 1,
			"cabin":     "ECONOMY",
		}
		resp, err := testApp.makeRequest("POST", "/api/v1/vouchers", voucherBody)
		if err != nil {
			t.Fatalf("Failed to create voucher: %v", err)
		}
		if resp.Code != http.StatusCreated {
			t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.Code)
		}

		var expiresAt *string
		err = testApp.DB.QueryRow("SELECT expires_at FROM vouchers WHERE code = ?", "VOUCHER_NO_EXPIRY").Scan(&expiresAt)
		if err != nil {
			t.Fatalf("Failed to query database: %v", err)
		}

		if expiresAt != nil {
			t.Errorf("Expected expires_at to be NULL, got: %v", *expiresAt)
		}
	})

	t.Run("Voucher with expires_at should have value in database", func(t *testing.T) {
		voucherBody := map[string]any{
			"code":       "VOUCHER_WITH_EXPIRY",
			"flight_id":  1,
			"cabin":      "ECONOMY",
			"expires_at": "2025-12-31T23:59:59Z",
		}
		resp, err := testApp.makeRequest("POST", "/api/v1/vouchers", voucherBody)
		if err != nil {
			t.Fatalf("Failed to create voucher: %v", err)
		}
		if resp.Code != http.StatusCreated {
			t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.Code)
		}

		var expiresAt *string
		err = testApp.DB.QueryRow("SELECT expires_at FROM vouchers WHERE code = ?", "VOUCHER_WITH_EXPIRY").Scan(&expiresAt)
		if err != nil {
			t.Fatalf("Failed to query database: %v", err)
		}

		if expiresAt == nil {
			t.Error("Expected expires_at to have value, got NULL")
		} else if *expiresAt != "2025-12-31T23:59:59Z" {
			t.Errorf("Expected expires_at to be '2025-12-31T23:59:59Z', got: %v", *expiresAt)
		}
	})
}
