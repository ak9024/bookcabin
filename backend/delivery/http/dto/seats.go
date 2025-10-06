package dto

type CreateBulkSeatRequest struct {
	FlightID int64    `json:"flight_id" validate:"required,gt=0"`
	Cabin    string   `json:"cabin" validate:"required,oneof=ECONOMY BUSINESS FIRST"`
	Labels   []string `json:"labels" validate:"required,min=1,dive,required"`
}
