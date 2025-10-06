package models

import "time"

type CreateBulkFlight struct {
	FlightNumbers []string  `json:"flight_numbers"` // e.g. ["GA133", "GA125"]
	DepDate       time.Time `json:"dep_date"`       // departure date
}

type Flight struct {
	ID       int64     `json:"id"`
	FlightNo string    `json:"flight_no"` // flight number, e.g. "GA133, GA125"
	DepDate  time.Time `json:"dep_date"`  // departure date
}

type Flights = []Flight
