package http

import (
	"backend/delivery/http/handler"

	"github.com/gofiber/fiber/v2"
)

func Routes(
	app *fiber.App,
	flightsHandler handler.FlightsHandler,
	seatsHandler handler.SeatsHandler,
	vouchersHandler handler.VouchersHandler,
) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// flights
	flights := v1.Group("/flights")
	flights.Post("/", flightsHandler.Create)
	flights.Get("/", flightsHandler.GetAll)

	// seats
	seats := v1.Group("/seats")
	seats.Get("/", seatsHandler.Get)
	seats.Post("/", seatsHandler.Create)

	// vouchers
	vouchers := v1.Group("/vouchers")
	vouchers.Post("/", vouchersHandler.Create)
	vouchers.Get("/", vouchersHandler.GetAll)
	vouchers.Post("/assigns", vouchersHandler.Assigns)
}
