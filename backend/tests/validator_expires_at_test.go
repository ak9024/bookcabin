package tests

import (
	"backend/delivery/http/dto"
	"backend/delivery/http/validator"
	"testing"
)

func TestValidateVoucherExpiresAt(t *testing.T) {
	validExpiresAt := "2025-12-31T23:59:59Z"
	invalidExpiresAt := "2025-12-31"

	tests := []struct {
		name        string
		request     dto.CreateNewVoucherRequest
		expectError bool
	}{
		{
			name: "Valid voucher with expires_at",
			request: dto.CreateNewVoucherRequest{
				Code:      "VOUCHER100",
				FlightID:  1,
				Cabin:     "ECONOMY",
				ExpiresAt: &validExpiresAt,
			},
			expectError: false,
		},
		{
			name: "Valid voucher without expires_at",
			request: dto.CreateNewVoucherRequest{
				Code:      "VOUCHER200",
				FlightID:  1,
				Cabin:     "ECONOMY",
				ExpiresAt: nil,
			},
			expectError: false,
		},
		{
			name: "Invalid expires_at format",
			request: dto.CreateNewVoucherRequest{
				Code:      "VOUCHER300",
				FlightID:  1,
				Cabin:     "ECONOMY",
				ExpiresAt: &invalidExpiresAt,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateStruct(&tt.request)
			if tt.expectError && err == nil {
				t.Errorf("Expected validation error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no validation error, got: %v", err)
			}
		})
	}
}
