package dto

type AssignVoucherRequest struct {
	VoucherCode string `json:"voucher_code"`
}

type CreateNewVoucherRequest struct {
	Code     string `json:"code"`
	FlightID int64  `json:"flight_id"`
	Cabin    string `json:"cabin"` // ECONOMY|BUSINESS|FIRST
}

type Voucher struct {
	ID         int64   `json:"id"`
	Code       string  `json:"code"`
	FlightID   int64   `json:"flight_id"`
	Cabin      string  `json:"cabin"`
	ExpiresAt  string  `json:"expires_at"` // voucher time periode
	Redeemed   int64   `json:"redeemed"`   // redeemed is used to flag or mark the voucher is used or not!
	RedeemedAt *string `json:"redeemed_at,omitempty"`
}

type Vouchers = []Voucher
