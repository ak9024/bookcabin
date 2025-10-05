package http

import (
	"backend/delivery/http/handler"

	"github.com/gofiber/fiber/v2"
)

func Routes(
	app *fiber.App,
	flightsHandler *handler.FlightsHandler,
	seatsHandler *handler.SeatsHandler,
	vouchersHandler *handler.VouchersHandler,
	assignmentsHandler *handler.AssignmentsHandler,
) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// flights
	flights := v1.Group("/flights")
	flights.Post("/", flightsHandler.Create)
	flights.Get("/", flightsHandler.Get)
}
