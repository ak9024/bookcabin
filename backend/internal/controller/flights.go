package controller

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type FlightsController interface {
	Create(ctx context.Context, flights *models.CreateBulkFlight) error
	Get(ctx context.Context) (models.Flights, error)
}

type flightsController struct {
	repo repository.FlightsRepository
}

func NewFlightsController(flightsRepo repository.FlightsRepository) FlightsController {
	return &flightsController{
		repo: flightsRepo,
	}
}

func (fc *flightsController) Create(ctx context.Context, flights *models.CreateBulkFlight) error {
	if err := fc.repo.Create(ctx, flights); err != nil {
		return err
	}

	return nil
}

func (fc *flightsController) Get(ctx context.Context) (models.Flights, error) {
	flights, err := fc.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return flights, nil
}
