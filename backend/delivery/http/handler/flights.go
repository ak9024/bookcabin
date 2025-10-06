package handler

import (
	"backend/delivery/http/dto"
	"backend/delivery/http/validator"
	"backend/internal/controller"
	"backend/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

type FlightsHandler interface {
	GetAll(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
}

type flightsHandler struct {
	fc controller.FlightsController
}

func NewFlightsHandler(flightsController controller.FlightsController) FlightsHandler {
	return &flightsHandler{fc: flightsController}
}

func (fh *flightsHandler) Create(c *fiber.Ctx) error {
	p := new(dto.CreateBulkFlightRequest)
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.JsonResponses{
			StatusCode: fiber.StatusBadRequest,
			Data:       err.Error(),
		})
	}

	if err := validator.ValidateStruct(p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.JsonResponses{
			StatusCode: fiber.StatusBadRequest,
			Data:       validator.FormatValidationErrors(err),
		})
	}

	depDate, err := time.Parse("2006-01-02", p.DepDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.JsonResponses{
			StatusCode: fiber.StatusBadRequest,
			Data:       "invalid date format, expected YYYY-MM-DD",
		})
	}

	if err := fh.fc.Create(c.Context(), &models.CreateBulkFlight{
		FlightNumbers: p.FlightNumbers,
		DepDate:       depDate,
	}); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.JsonResponses{
			StatusCode: fiber.StatusBadRequest,
			Data:       err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.JsonResponses{
		StatusCode: fiber.StatusCreated,
		Data:       "success to create a new flights",
	})
}

func (fh *flightsHandler) GetAll(c *fiber.Ctx) error {
	flights, err := fh.fc.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(dto.JsonResponses{
		StatusCode: fiber.StatusOK,
		Data:       flights,
	})
}
