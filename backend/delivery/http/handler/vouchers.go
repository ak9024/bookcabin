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
	p := new(dto.CreateNewVoucherRequest)
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.JsonResponses{
			StatusCode: fiber.StatusBadRequest,
			Data:       err.Error(),
		})
	}

	if err := vh.vc.Create(c.Context(), &models.CreateNewVoucher{
		Code:     p.Code,
		FlightID: p.FlightID,
		Cabin:    p.Cabin,
	}); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.JsonResponses{
			StatusCode: fiber.StatusBadRequest,
			Data:       err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.JsonResponses{
		StatusCode: fiber.StatusCreated,
		Data:       "success to create a new voucher",
	})
}

func (vh *vouchersHandler) Assigns(c *fiber.Ctx) error {
	p := new(dto.AssignVoucherRequest)
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.JsonResponses{
			StatusCode: fiber.StatusBadRequest,
			Data:       err.Error(),
		})
	}

	voucher, err := vh.vc.Assigns(c.Context(), &models.AssignsRandomVoucher{
		VoucherCode: p.VoucherCode,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.JsonResponses{
			StatusCode: fiber.StatusBadRequest,
			Data:       err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.JsonResponses{
		StatusCode: fiber.StatusCreated,
		Data:       voucher,
	})
}

func (vh *vouchersHandler) GetAll(c *fiber.Ctx) error {
	rows, err := vh.vc.GetAll(c.Context())

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.JsonResponses{
			StatusCode: fiber.StatusBadRequest,
			Data:       err.Error(),
		})
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

	return c.Status(fiber.StatusOK).JSON(dto.JsonResponses{
		StatusCode: fiber.StatusOK,
		Data:       vouchers,
	})
}
