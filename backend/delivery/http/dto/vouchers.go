package dto

type AssignVoucherPayload struct {
	VoucherCode string `json:"voucher_code"`
}

type CreateNewVoucherPayload struct {
	Code     string `json:"code"`
	FlightID int64  `json:"flight_id"`
	Cabin    string `json:"cabin"` // ECONOMY|BUSINESS|FIRST
}
