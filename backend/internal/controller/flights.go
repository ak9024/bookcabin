package controller

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type FlightsController interface {
	Create(ctx context.Context, flights *models.CreateBulkFlight) error
	GetAll(ctx context.Context) (models.Flights, error)
}

type flightsController struct {
	fr repository.FlightsRepository
}

func NewFlightsController(fr repository.FlightsRepository) FlightsController {
	return &flightsController{
		fr: fr,
	}
}

func (fc *flightsController) Create(ctx context.Context, flights *models.CreateBulkFlight) error {
	if err := fc.fr.Create(ctx, flights); err != nil {
		return err
	}

	return nil
}

func (fc *flightsController) GetAll(ctx context.Context) (models.Flights, error) {
	flights, err := fc.fr.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return flights, nil
}
