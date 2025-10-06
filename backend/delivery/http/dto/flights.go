package dto

import "time"

type CreateBulkFlightRequest struct {
	FlightNumbers []string  `json:"flight_numbers"` // e.g. ["GA133", "GA125"]
	DepDate       time.Time `json:"dep_date"`       // departure date
}

func (cbfr *CreateBulkFlightRequest) Validate() bool {
	if len(cbfr.FlightNumbers) == 0 || cbfr.DepDate.String() == "" {
		return false
	} else {
		return true
	}
}
