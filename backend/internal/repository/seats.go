package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"errors"
	"strings"
)

type SeatRepository interface {
	Create(ctx context.Context, cbs *models.CreateBulkSeat) error
	GetAll(ctx context.Context) (*models.Seats, error)
}

type seatRepository struct {
	db *sql.DB
}

func NewSeatRepository(db *sql.DB) SeatRepository {
	return &seatRepository{
		db: db,
	}
}

func (sr *seatRepository) flightExists(ctx context.Context, flightID int64) (bool, error) {
	var exists bool
	err := sr.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM flights WHERE id = ?)", flightID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (sr *seatRepository) Create(ctx context.Context, cbs *models.CreateBulkSeat) error {
	// Validate flight exists
	exists, err := sr.flightExists(ctx, cbs.FlightID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("flight not found")
	}

	tx, _ := sr.db.Begin()
	defer tx.Rollback()

	for _, l := range cbs.Labels {
		l = strings.ToUpper(strings.TrimSpace(l))
		if _, err := tx.Exec(`INSERT INTO seats(flight_id, label, cabin) VALUES(?,?,?)`, cbs.FlightID, l, cbs.Cabin); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (sr *seatRepository) GetAll(ctx context.Context) (*models.Seats, error) {
	rows, err := sr.db.Query("SELECT id, flight_id, label, cabin FROM seats")
	if err != nil {
		return nil, err
	}

	var seats models.Seats
	for rows.Next() {
		var seat models.Seat
		if err := rows.Scan(&seat.ID, &seat.FlightID, &seat.Label, &seat.Cabin); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}

	return &seats, nil
}
