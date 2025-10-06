package handler

import "github.com/gofiber/fiber/v2"

type SeatsHandler interface {
	Get(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
}

type seatsHandler struct{}

func NewSeatsHandler() SeatsHandler {
	return &seatsHandler{}
}

func (sh *seatsHandler) Create(c *fiber.Ctx) error {
	return nil
}

func (sh *seatsHandler) Get(c *fiber.Ctx) error {
	return nil
}
