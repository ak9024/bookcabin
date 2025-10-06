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
	repo repository.VouchersRepository
}

func NewVouchersController(vouchersRepository repository.VouchersRepository) VouchersController {
	return &vouchersController{
		repo: vouchersRepository,
	}
}

func (vc *vouchersController) Create(ctx context.Context, cnv *models.CreateNewVoucher) error {
	if err := vc.repo.Create(ctx, cnv); err != nil {
		return err
	}

	return nil
}

func (vc *vouchersController) Assigns(ctx context.Context, arv *models.AssignsRandomVoucher) (*models.VoucherAssigment, error) {
	voucher, err := vc.repo.Assigns(ctx, arv)
	if err != nil {
		return nil, err
	}

	return voucher, nil
}

func (vc *vouchersController) GetAll(ctx context.Context) (*models.Vouchers, error) {
	vouchers, err := vc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return vouchers, nil
}
