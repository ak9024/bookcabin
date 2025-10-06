package dto

type CreateBulkFlightRequest struct {
	FlightNumbers []string `json:"flight_numbers"` // e.g. ["GA133", "GA125"]
	DepDate       string   `json:"dep_date"`       // departure date in YYYY-MM-DD format
}

func (cbfr *CreateBulkFlightRequest) Validate() bool {
	if len(cbfr.FlightNumbers) == 0 || cbfr.DepDate == "" {
		return false
	} else {
		return true
	}
}
