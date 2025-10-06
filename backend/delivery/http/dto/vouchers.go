package dto

type AssignVoucherRequest struct {
	VoucherCode string `json:"voucher_code" validate:"required,min=1"`
}

type CreateNewVoucherRequest struct {
	Code      string  `json:"code" validate:"required,min=1"`
	FlightID  int64   `json:"flight_id" validate:"required,gt=0"`
	Cabin     string  `json:"cabin" validate:"required,oneof=ECONOMY BUSINESS FIRST"` // ECONOMY|BUSINESS|FIRST
	ExpiresAt *string `json:"expires_at,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type Voucher struct {
	ID         int64   `json:"id"`
	Code       string  `json:"code"`
	FlightID   int64   `json:"flight_id"`
	Cabin      string  `json:"cabin"`
	ExpiresAt  *string `json:"expires_at,omitempty"` // voucher time periode
	Redeemed   int64   `json:"redeemed"`             // redeemed is used to flag or mark the voucher is used or not!
	RedeemedAt *string `json:"redeemed_at,omitempty"`
}

type Vouchers = []Voucher
