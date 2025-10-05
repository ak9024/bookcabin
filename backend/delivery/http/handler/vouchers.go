package handler

import (
	"backend/delivery/http/dto"
	"backend/internal/controller"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type VouchersHandler struct {
	vc controller.VouchersController
}

func NewVouchersHandler(vouchersController controller.VouchersController) *VouchersHandler {
	return &VouchersHandler{
		vc: vouchersController,
	}
}

func (vh *VouchersHandler) Create(c *fiber.Ctx) error {
	vp := new(dto.CreateNewVoucherPayload)
	if err := c.BodyParser(&vp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	if err := vh.vc.Create(c.Context(), &models.CreateNewVoucher{
		Code:     vp.Code,
		FlightID: vp.FlightID,
		Cabin:    vp.Cabin,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (vh *VouchersHandler) Assign(c *fiber.Ctx) error {
	vp := new(dto.AssignVoucherPayload)
	if err := c.BodyParser(&vp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	voucher, err := vh.vc.Assign(c.Context(), &models.AssignRandomVoucher{
		VoucherCode: vp.VoucherCode,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(voucher)
}
