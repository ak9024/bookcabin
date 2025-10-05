package models

import (
	"database/sql"
)

type (
	Voucher struct {
		ID         int64          `json:"id"`
		Code       string         `json:"code"`
		FlightID   int64          `json:"flight_id"`
		Cabin      string         `json:"cabin"`
		ExpiresAt  sql.NullString `json:"expires_at"` // voucher time periode
		Redeemed   int64          `json:"redeemed"`   // redeemed is used to flag or mark the voucher is used or not!
		RedeemedAt *string        `json:"redeemed_at,omitempty"`
	}

	CreateNewVoucher struct {
		Code     string `json:"code"`
		FlightID int64  `json:"flight_id"`
		Cabin    string `json:"cabin"` // ECONOMY|BUSINESS|FIRST
	}

	AssignRandomVoucher struct {
		VoucherCode string `json:"voucher_code"`
	}
)

type Vouchers = []Voucher
