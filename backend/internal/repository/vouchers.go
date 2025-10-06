package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

type VouchersRepository interface {
	Assigns(ctx context.Context, arv *models.AssignsRandomVoucher) (*models.VoucherAssigment, error)
	Create(ctx context.Context, cnv *models.CreateNewVoucher) error
	GetAll(ctx context.Context) (*models.Vouchers, error)
}

type vouchersRepository struct {
	db *sql.DB
}

func NewVouchersRepository(db *sql.DB) VouchersRepository {
	return &vouchersRepository{
		db,
	}
}

func (vr *vouchersRepository) Create(ctx context.Context, cnv *models.CreateNewVoucher) error {
	var seatID int64
	err := vr.db.QueryRow(`SELECT id FROM seats WHERE flight_id=? AND cabin=? LIMIT 1`, cnv.FlightID, cnv.Cabin).Scan(&seatID)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("no seats available for this flight and cabin")
	} else if err != nil {
		return err
	}

	if _, err := vr.db.Exec(`INSERT INTO vouchers(code, flight_id, cabin) VALUES(?, ?, ?)`, cnv.Code, cnv.FlightID, cnv.Cabin); err != nil {
		return err
	}

	return nil
}

func (vr *vouchersRepository) Assigns(ctx context.Context, arv *models.AssignsRandomVoucher) (*models.VoucherAssigment, error) {
	const maxAttempts = 3
	var lastError error
	var result models.VoucherAssigment

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		tx, err := vr.db.BeginTx(ctx, &sql.TxOptions{})
		if err != nil {
			return nil, err
		}
		defer tx.Rollback()

		var v models.Voucher
		err = tx.QueryRowContext(ctx, `SELECT id, flight_id, cabin, redeemed, COALESCE(expires_at,'')
		FROM vouchers WHERE code=?`, arv.VoucherCode).Scan(&v.ID, &v.FlightID, &v.Cabin, &v.Redeemed, &v.ExpiresAt)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("voucher not found!")
		} else if err != nil {
			return nil, err
		}

		if v.Redeemed == 1 {
			return nil, errors.New("voucher already redeemed!")
		}

		if v.ExpiresAt.Valid && v.ExpiresAt.String != "" {
			if t, e := time.Parse(time.RFC3339, v.ExpiresAt.String); e == nil && time.Now().After(t) {
				return nil, errors.New("voucher expired!")
			}
		}

		candidateSeatID, err := pickSeatRandomly(ctx, tx, v.FlightID, v.Cabin)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				lastError = errors.New("no available seats in cabin!")
				continue
			}
		}

		if _, err := tx.ExecContext(ctx,
			`INSERT INTO seat_assignments(voucher_id, seat_id) VALUES(?, ?)
				 ON CONFLICT(seat_id) DO NOTHING`, v.ID, candidateSeatID); err != nil {
			lastError = err
			continue
		}

		var count int
		if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM seat_assignments WHERE voucher_id=?`, v.ID).Scan(&count); err != nil {
			lastError = err
			continue
		}

		if count == 0 {
			lastError = errors.New("seat taken concurrently")
			continue
		}

		if _, err := tx.ExecContext(ctx, `UPDATE seats SET is_assigned=1 WHERE id=?`, candidateSeatID); err != nil {
			lastError = err
			continue
		}

		if _, err := tx.ExecContext(ctx, `UPDATE vouchers SET redeemed=1, redeemed_at=datetime('now') WHERE id=?`, v.ID); err != nil {
			lastError = err
			continue
		}

		if err := tx.Commit(); err != nil {
			return nil, err
		}

		var seatLabel string
		_ = vr.db.QueryRow(`SELECT label FROM seats WHERE id=?`, candidateSeatID).Scan(&seatLabel)

		result := &models.VoucherAssigment{
			VoucherCode: arv.VoucherCode,
			Cabin:       v.Cabin,
			SeatID:      candidateSeatID,
			SeatLabel:   seatLabel,
		}

		return result, nil
	}

	if lastError != nil {
		return nil, lastError
	}

	return &result, nil
}

func pickSeatRandomly(ctx context.Context, tx *sql.Tx, flightID int64, cabin string) (int64, error) {
	var seatID int64
	err := tx.QueryRowContext(ctx, `SELECT id FROM seats WHERE flight_id=? AND cabin=? AND is_assigned=0 ORDER BY random() LIMIT 1`, flightID, cabin).Scan(&seatID)
	if err != nil {
		return 0, err
	}
	return seatID, nil
}

func (vr *vouchersRepository) GetAll(ctx context.Context) (*models.Vouchers, error) {
	rows, err := vr.db.Query("SELECT id, flight_id, code, cabin, redeemed, expires_at, redeemed_at FROM vouchers")
	if err != nil {
		return nil, err
	}

	var vouchers models.Vouchers
	for rows.Next() {
		var voucher models.Voucher
		if err := rows.Scan(&voucher.ID, &voucher.FlightID, &voucher.Code, &voucher.Cabin, &voucher.Redeemed, &voucher.ExpiresAt, &voucher.RedeemedAt); err != nil {
			return nil, err
		}

		vouchers = append(vouchers, voucher)
	}

	return &vouchers, nil
}
