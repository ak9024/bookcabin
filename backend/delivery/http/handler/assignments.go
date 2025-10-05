package handler

import "github.com/gofiber/fiber/v2"

type AssignmentsHandler struct{}

func NewAssignmentsHandler() *AssignmentsHandler {
	return &AssignmentsHandler{}
}

func (ah *AssignmentsHandler) Create(c *fiber.Ctx) error {
	return nil
}
