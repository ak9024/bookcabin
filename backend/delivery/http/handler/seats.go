package handler

import (
	"backend/delivery/http/dto"
	"backend/internal/controller"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type SeatsHandler interface {
	GetAll(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
}

type seatsHandler struct {
	sc controller.SeatController
}

func NewSeatsHandler(sc controller.SeatController) SeatsHandler {
	return &seatsHandler{
		sc: sc,
	}
}

func (sh *seatsHandler) Create(c *fiber.Ctx) error {
	p := new(dto.CreateBulkSeatRequest)
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	if err := sh.sc.Create(c.Context(), &models.CreateBulkSeat{
		FlightID: p.FlightID,
		Cabin:    p.Cabin,
		Labels:   p.Labels,
	}); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (sh *seatsHandler) GetAll(c *fiber.Ctx) error {
	seats, err := sh.sc.GetAll(c.Context())

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(seats)
}
