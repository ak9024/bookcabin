package controller

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type VouchersController interface {
	Assign(ctx context.Context, arv *models.AssignRandomVoucher) (*models.Voucher, error)
	Create(ctx context.Context, cnv *models.CreateNewVoucher) error
}

type vouchersController struct {
	repo repository.VouchersRepository
}

func NewVouchersController(vouchersRepository repository.VouchersRepository) VouchersController {
	return &vouchersController{
		repo: vouchersRepository,
	}
}

func (vc *vouchersController) Assign(ctx context.Context, arv *models.AssignRandomVoucher) (*models.Voucher, error) {
	const maxAttempts = 3
	var lastError error
	var v models.Voucher

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		voucher, err := vc.repo.Assign(ctx, arv)
		if err != nil {
			lastError = err
			continue
		}

		v = *voucher
	}

	if lastError != nil {
		return nil, lastError
	}

	return &v, nil
}

func (vc *vouchersController) Create(ctx context.Context, cnv *models.CreateNewVoucher) error {
	if err := vc.repo.Create(ctx, cnv); err != nil {
		return err
	}

	return nil
}
