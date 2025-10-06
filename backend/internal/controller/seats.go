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
	sr repository.SeatRepository
}

func NewSeatController(sr repository.SeatRepository) SeatController {
	return &seatController{
		sr: sr,
	}
}

func (sc *seatController) Create(ctx context.Context, cbs *models.CreateBulkSeat) error {
	if err := sc.sr.Create(ctx, cbs); err != nil {
		return err
	}

	return nil
}

func (sc *seatController) GetAll(ctx context.Context) (*models.Seats, error) {
	seats, err := sc.sr.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return seats, nil
}
