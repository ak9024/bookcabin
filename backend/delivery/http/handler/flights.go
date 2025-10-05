package handler

import (
	"backend/delivery/http/dto"
	"backend/internal/controller"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type FlightsHandler struct {
	fc controller.FlightsController
}

func NewFlightsHandler(flightsController controller.FlightsController) *FlightsHandler {
	return &FlightsHandler{fc: flightsController}
}

func (fh *FlightsHandler) Create(c *fiber.Ctx) error {
	f := new(dto.FlightPayload)
	if err := c.BodyParser(&f); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	if !f.Validate() {
		return c.Status(fiber.StatusBadRequest).JSON("invalid json")
	}

	if err := fh.fc.Create(c.Context(), &models.CreateBulkFlight{
		FlightNumbers: f.FlightNumbers,
		DepDate:       f.DepDate,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (fh *FlightsHandler) Get(c *fiber.Ctx) error {
	flights, err := fh.fc.Get(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(flights)
}
