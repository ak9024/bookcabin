package controller

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type VouchersController interface {
	Create(ctx context.Context, cnv *models.CreateNewVoucher) error
	Assigns(ctx context.Context, arv *models.AssignsRandomVoucher) (*models.VoucherAssigment, error)
	GetAll(ctx context.Context) (*models.Vouchers, error)
}

type vouchersController struct {
	vr repository.VouchersRepository
}

func NewVouchersController(vr repository.VouchersRepository) VouchersController {
	return &vouchersController{
		vr: vr,
	}
}

func (vc *vouchersController) Create(ctx context.Context, cnv *models.CreateNewVoucher) error {
	if err := vc.vr.Create(ctx, cnv); err != nil {
		return err
	}

	return nil
}

func (vc *vouchersController) Assigns(ctx context.Context, arv *models.AssignsRandomVoucher) (*models.VoucherAssigment, error) {
	voucher, err := vc.vr.Assigns(ctx, arv)
	if err != nil {
		return nil, err
	}

	return voucher, nil
}

func (vc *vouchersController) GetAll(ctx context.Context) (*models.Vouchers, error) {
	vouchers, err := vc.vr.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return vouchers, nil
}
