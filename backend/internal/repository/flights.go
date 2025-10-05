package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"strings"
	"time"
)

type FlightsRepository interface {
	Create(ctx context.Context, flight *models.CreateBulkFlight) error
	Get(ctx context.Context) (models.Flights, error)
}

type flightsRepository struct {
	db *sql.DB
}

func NewFlightsRepository(db *sql.DB) FlightsRepository {
	return &flightsRepository{
		db,
	}
}

func (fr *flightsRepository) Create(ctx context.Context, flight *models.CreateBulkFlight) error {
	tx, _ := fr.db.Begin()
	defer tx.Rollback()

	for _, fn := range flight.FlightNumbers {
		fn = strings.ToUpper(strings.TrimSpace(fn))
		if _, err := tx.Exec(`INSERT INTO flights(flight_no, dep_date) VALUES(?,?)`, fn, flight.DepDate.Format(time.RFC3339)); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (fr *flightsRepository) Get(ctx context.Context) (models.Flights, error) {
	rows, err := fr.db.Query("SELECT id, flight_no, dep_date FROM flights")
	if err != nil {
		return nil, err
	}

	var flights models.Flights

	for rows.Next() {
		var flight models.Flight
		var depDateStr string

		if err := rows.Scan(&flight.ID, &flight.FlightNo, &depDateStr); err != nil {
			return nil, err
		}

		if parsedTime, err := time.Parse(time.RFC3339, depDateStr); err == nil {
			flight.DepDate = parsedTime
		}

		flights = append(flights, flight)
	}

	return flights, nil
}
