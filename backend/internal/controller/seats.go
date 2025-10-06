package controller

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type SeatController interface {
	Create(ctx context.Context, cbs *models.CreateBulkSeat) error
	GetAll(ctx context.Context) (*models.Seats, error)
}

type seatController struct {
	seatRepo repository.SeatRepository
}

func NewSeatController(seatRepo repository.SeatRepository) SeatController {
	return &seatController{
		seatRepo: seatRepo,
	}
}

func (sc *seatController) Create(ctx context.Context, cbs *models.CreateBulkSeat) error {
	if err := sc.seatRepo.Create(ctx, cbs); err != nil {
		return err
	}

	return nil
}

func (sc *seatController) GetAll(ctx context.Context) (*models.Seats, error) {
	seats, err := sc.seatRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return seats, nil
}
