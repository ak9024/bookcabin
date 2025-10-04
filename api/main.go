package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func initDB() {
	var err error

	DB, err = sql.Open("sqlite3", "./bookcabin.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS flights(
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  flight_no  TEXT NOT NULL,
  dep_date   TEXT NOT NULL, -- YYYY-MM-DD
  UNIQUE(flight_no, dep_date)
);

CREATE TABLE IF NOT EXISTS seats(
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  flight_id    INTEGER NOT NULL REFERENCES flights(id) ON DELETE CASCADE,
  label        TEXT NOT NULL, -- e.g., 12A
  cabin        TEXT NOT NULL CHECK (cabin IN ('ECONOMY','BUSINESS','FIRST')),
  is_assigned  INTEGER NOT NULL DEFAULT 0,
  UNIQUE(flight_id, label)
);

CREATE TABLE IF NOT EXISTS vouchers(
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  code         TEXT NOT NULL UNIQUE,
  flight_id    INTEGER NOT NULL REFERENCES flights(id) ON DELETE CASCADE,
  cabin        TEXT NOT NULL CHECK (cabin IN ('ECONOMY','BUSINESS','FIRST')),
  redeemed     INTEGER NOT NULL DEFAULT 0,
  expires_at   TEXT,   -- RFC3339
  redeemed_at  TEXT
);

CREATE TABLE IF NOT EXISTS seat_assignments(
  voucher_id   INTEGER NOT NULL REFERENCES vouchers(id) ON DELETE CASCADE,
  seat_id      INTEGER NOT NULL REFERENCES seats(id) ON DELETE CASCADE,
  assigned_at  TEXT NOT NULL DEFAULT (datetime('now')),
  PRIMARY KEY (voucher_id),
  UNIQUE (seat_id)
);

CREATE INDEX IF NOT EXISTS idx_seats_flight_cabin ON seats(flight_id, cabin);`

	_, err = DB.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}

type (
	Flight struct {
		ID       int64  `json:"id"`
		FlightNo string `json:"flight_no"` // flight number
		DepDate  string `json:"dep_date"`  // departure date
	}

	Seat struct {
		ID         int64  `json:"id"`
		FlightID   int64  `json:"flight_id"`
		Label      string `json:"label"` // e.g., "12A"
		Cabin      string `json:"cabin"` // ECONOMY|BUSINESS|FIRST
		IsAssigned bool   `json:"is_assigned"`
	}

	Vourcher struct {
		ID          int64   `json:"id"`
		Code        string  `json:"code"`
		FlightID    int64   `json:"flight_id"`
		Cabin       string  `json:"cabin"`
		ExpiresAt   *string `json:"expires_at,omitempty"`  // voucher time periode
		Redeemed    bool    `json:"redeemed"`              // redeemed is used to flag or mark the voucher is used or not!
		ReadeemedAt *string `json:"redeemed_at,omitempty"` // fill with date RFC3339 when the voucher is used
	}

	Assigment struct {
		ID          int64  `json:"id"`
		VoucherCode string `json:"voucher_code"`
		SeatID      int64  `json:"seat_id"`
		SeatLabel   string `json:"seat_label"`
		Cabin       string `json:"cabin"`
		Status      bool   `json:"status"`
	}
)

func getAllFlights(c *fiber.Ctx) error {
	flights := []Flight{
		{},
	}
	return c.JSON(flights)
}

func getAllSeats(c *fiber.Ctx) error {
	seats := []Seat{
		{},
	}
	return c.JSON(seats)
}

func getAllVouchers(c *fiber.Ctx) error {
	vouchers := []Vourcher{
		{},
	}
	return c.JSON(vouchers)
}

func getAllAssigments(c *fiber.Ctx) error {
	assigments := []Assigment{
		{},
	}
	return c.JSON(assigments)
}

func main() {
	initDB()
	defer DB.Close()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(healthcheck.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// seats
	seats := v1.Group("/seats")
	seats.Get("/", getAllSeats)

	// vouchers
	vouchers := v1.Group("/vouchers")
	vouchers.Get("/", getAllVouchers)

	// flights
	flights := v1.Group("/flights")
	flights.Get("/", getAllFlights)

	// assignments
	assigments := v1.Group("/assigments")
	assigments.Get("/", getAllAssigments)

	app.Listen(":8080")
}
