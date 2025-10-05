package handler

import "github.com/gofiber/fiber/v2"

type VouchersHandler struct{}

func NewVouchersHandler() *VouchersHandler {
	return &VouchersHandler{}
}

func (vh *VouchersHandler) Create(c *fiber.Ctx) error {
	return nil
}
