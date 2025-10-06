package dto

type CreateBulkFlightRequest struct {
	FlightNumbers []string `json:"flight_numbers" validate:"required,min=1,dive,required"` // e.g. ["GA133", "GA125"]
	DepDate       string   `json:"dep_date" validate:"required,datetime=2006-01-02"`       // departure date in YYYY-MM-DD format
}
