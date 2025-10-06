package models

type Seat struct {
	ID       int64  `json:"id"`
	FlightID int64  `json:"flight_id"`
	Label    string `json:"label"`
	Cabin    string `json:"cabin"`
}

type Seats = []Seat

type CreateBulkSeat struct {
	FlightID int64    `json:"flight_id"`
	Cabin    string   `json:"cabin"`
	Labels   []string `json:"labels"`
}
