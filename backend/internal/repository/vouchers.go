package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

type VouchersRepository interface {
	Assign(ctx context.Context, arv *models.AssignRandomVoucher) (*models.Voucher, error)
	Create(ctx context.Context, cnv *models.CreateNewVoucher) error
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
	if _, err := vr.db.Exec(`INSERT INTO vouchers(code, flight_id, cabin) VALUES(?, ?, ?)`, cnv.Code, cnv.FlightID, cnv.Cabin); err != nil {
		return err
	}

	return nil
}

func (vr *vouchersRepository) Assign(ctx context.Context, arv *models.AssignRandomVoucher) (*models.Voucher, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
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
			return nil, errors.New("No available seats in cabin!")
		}
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO seat_assignments(voucher_id, seat_id) VALUES(?, ?)
				 ON CONFLICT(seat_id) DO NOTHING`, v.ID, candidateSeatID); err != nil {
	}

	var count int
	if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM seat_assignments WHERE voucher_id=?`, v.ID).Scan(&count); err != nil {
	}

	if count == 0 {
		return nil, errors.New("seat taken concurrently!")
	}

	if _, err := tx.ExecContext(ctx, `UPDATE seats SET is_assigned=1 WHERE id=?`, candidateSeatID); err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `UPDATE vouchers SET redeemed=1, redeemed_at=datetime('now') WHERE id=?`, v.ID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var seatLabel string
	_ = vr.db.QueryRow(`SELECT label FROM seats WHERE id=?`, candidateSeatID).Scan(&seatLabel)

	return &v, nil
}

func pickSeatRandomly(ctx context.Context, tx *sql.Tx, flightID int64, cabin string) (int64, error) {
	var seatID int64
	err := tx.QueryRowContext(ctx, `SELECT id FROM seats WHERE flight_id=? AND cabin=? AND is_assigned=0 ORDER BY random() LIMIT 1`, flightID, cabin).Scan(&seatID)
	if err != nil {
		return 0, err
	}
	return seatID, nil
}
