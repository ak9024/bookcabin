package dto

type CreateBulkSeatRequest struct {
	FlightID int64    `json:"flight_id"`
	Cabin    string   `json:"cabin"`
	Labels   []string `json:"labels"`
}
