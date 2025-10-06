package handler

import (
	"backend/delivery/http/dto"
	"backend/internal/controller"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type VouchersHandler interface {
	Create(c *fiber.Ctx) error
	Assigns(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}

type vouchersHandler struct {
	vc controller.VouchersController
}

func NewVouchersHandler(vouchersController controller.VouchersController) VouchersHandler {
	return &vouchersHandler{
		vc: vouchersController,
	}
}

func (vh *vouchersHandler) Create(c *fiber.Ctx) error {
	vp := new(dto.CreateNewVoucherPayload)
	if err := c.BodyParser(&vp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	if err := vh.vc.Create(c.Context(), &models.CreateNewVoucher{
		Code:     vp.Code,
		FlightID: vp.FlightID,
		Cabin:    vp.Cabin,
	}); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (vh *vouchersHandler) Assigns(c *fiber.Ctx) error {
	vp := new(dto.AssignVoucherPayload)
	if err := c.BodyParser(&vp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	rows, err := vh.vc.Assigns(c.Context(), &models.AssignsRandomVoucher{
		VoucherCode: vp.VoucherCode,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(rows)
}

func (vh *vouchersHandler) GetAll(c *fiber.Ctx) error {
	rows, err := vh.vc.GetAll(c.Context())

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	var vouchers dto.Vouchers
	if rows != nil {
		for _, v := range *rows {
			voucher := dto.Voucher{
				ID:         v.ID,
				FlightID:   v.FlightID,
				Code:       v.Code,
				Cabin:      v.Cabin,
				Redeemed:   v.Redeemed,
				RedeemedAt: v.RedeemedAt,
			}

			if v.ExpiresAt.Valid {
				voucher.ExpiresAt = v.ExpiresAt.String
			}

			vouchers = append(vouchers, voucher)
		}
	}

	return c.Status(fiber.StatusOK).JSON(vouchers)
}
