package dto

import "time"

type FlightPayload struct {
	FlightNumbers []string  `json:"flight_numbers"` // e.g. ["GA133", "GA125"]
	DepDate       time.Time `json:"dep_date"`       // departure date
}

func (f *FlightPayload) Validate() bool {
	if len(f.FlightNumbers) == 0 || f.DepDate.String() == "" {
		return false
	} else {
		return true
	}
}
