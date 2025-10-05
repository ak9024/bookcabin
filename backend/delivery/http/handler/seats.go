package handler

import "github.com/gofiber/fiber/v2"

type SeatsHandler struct{}

func NewSeatsHandler() *SeatsHandler {
	return &SeatsHandler{}
}

func (sh *SeatsHandler) Create(c *fiber.Ctx) error {
	return nil
}
