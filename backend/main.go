package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Config struct {
	Port   string
	DBPath string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./bookcabin.db"
	}

	return &Config{
		Port:   port,
		DBPath: dbPath,
	}
}

func initDB(dbPath string) {
	var err error

	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS flights(
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  flight_no  TEXT NOT NULL,
  dep_date   TEXT NOT NULL, -- RFC3339
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
  assigned_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),  -- RFC3339
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
		ID       int64     `json:"id"`
		FlightNo string    `json:"flight_no"` // flight number, e.g. "GA133, GA125"
		DepDate  time.Time `json:"dep_date"`  // departure date
	}

	FlightPayload struct {
		FlightNumbers []string  `json:"flight_numbers"` // e.g. ["GA133", "GA125"]
		DepDate       time.Time `json:"dep_date"`       // departure date
	}

	Seat struct {
		ID         int64  `json:"id"`
		FlightID   int64  `json:"flight_id"`
		Label      string `json:"label"` // e.g., "12A"
		Cabin      string `json:"cabin"` // ECONOMY|BUSINESS|FIRST
		IsAssigned bool   `json:"is_assigned"`
	}

	SeatPayload struct {
		FlightID int64    `json:"flight_id"`
		Cabin    string   `json:"cabin"`  // ECONOMY|BUSINESS|FIRST
		Labels   []string `json:"labels"` // ["1A", "1B", ...]
	}

	Voucher struct {
		ID          int64      `json:"id"`
		Code        string     `json:"code"`
		FlightID    int64      `json:"flight_id"`
		Cabin       string     `json:"cabin"`
		ExpiresAt   *time.Time `json:"expires_at,omitempty"`  // voucher time periode (RFC3339)
		Redeemed    bool       `json:"redeemed"`              // redeemed is used to flag or mark the voucher is used or not!
		ReadeemedAt *time.Time `json:"redeemed_at,omitempty"` // fill with date RFC3339 when the voucher is used
	}

	VoucherPayload struct {
		Code      string    `json:"code"`
		FlightID  int64     `json:"flight_id"`
		Cabin     string    `json:"cabin"` // ECONOMY|BUSINESS|FIRST
		ExpiresAt time.Time `json:"expires_at"`
	}

	Assigment struct {
		ID          int64  `json:"id"`
		VoucherCode string `json:"voucher_code"`
		SeatID      int64  `json:"seat_id"`
		SeatLabel   string `json:"seat_label"`
		Cabin       string `json:"cabin"`
		Status      bool   `json:"status"`
	}

	ErrorResponse struct {
		Code int32 `json:"code"`
		Data any   `json:"data"`
	}

	SucccessResponse struct {
		Code int32 `json:"code"`
		Data any   `json:"data"`
	}
)

func (sp *SeatPayload) validate() bool {
	if sp.FlightID == 0 || len(sp.Labels) == 0 || (sp.Cabin != "BUSINESS" && sp.Cabin != "ECONOMY" && sp.Cabin != "FIRST") {
		return false
	} else {
		return true
	}
}

func (fp *FlightPayload) validate() bool {
	if len(fp.FlightNumbers) == 0 {
		return false
	} else {
		return true
	}
}

func (vp *VoucherPayload) validate() bool {
	if vp.Code == "" || (vp.Cabin != "BUSINESS" && vp.Cabin != "ECONOMY" && vp.Cabin != "FIRST") || vp.FlightID == 0 {
		return false
	} else {
		return true
	}
}

func getAllFlights(c *fiber.Ctx) error {
	rows, err := DB.Query("SELECT id, flight_no, dep_date FROM flights")
	if err != nil {
		var er ErrorResponse
		er.Code = fiber.StatusInternalServerError
		er.Data = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(er)
	}

	var flights []Flight
	for rows.Next() {
		var flight Flight
		var depDateStr string

		if err := rows.Scan(&flight.ID, &flight.FlightNo, &depDateStr); err != nil {
			var er ErrorResponse
			er.Code = fiber.StatusInternalServerError
			er.Data = err.Error()
			return c.Status(fiber.StatusInternalServerError).JSON(er)
		}

		if parsedTime, err := time.Parse(time.RFC3339, depDateStr); err == nil {
			flight.DepDate = parsedTime
		}

		flights = append(flights, flight)
	}

	var sr SucccessResponse
	sr.Code = fiber.StatusOK
	sr.Data = flights
	return c.Status(fiber.StatusOK).JSON(sr)
}

func createNewFlight(c *fiber.Ctx) error {
	fp := new(FlightPayload)

	if err := c.BodyParser(&fp); err != nil {
		var er ErrorResponse
		er.Code = fiber.StatusInternalServerError
		er.Data = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(er)
	}

	fp.DepDate = fp.DepDate.UTC()

	if !fp.validate() {
		var er ErrorResponse
		er.Code = fiber.StatusBadRequest
		er.Data = "Invalid JSON"
		return c.Status(fiber.StatusBadRequest).JSON(er)
	} else {
		tx, _ := DB.Begin()
		defer tx.Rollback()

		for _, fn := range fp.FlightNumbers {
			fn = strings.ToUpper(strings.TrimSpace(fn))
			if _, err := tx.Exec(`INSERT INTO flights(flight_no, dep_date) VALUES(?,?)`, fn, fp.DepDate.Format(time.RFC3339)); err != nil {
				var er ErrorResponse
				er.Code = fiber.StatusInternalServerError
				er.Data = err.Error()
				return c.Status(fiber.StatusInternalServerError).JSON(er)
			}
		}

		if err := tx.Commit(); err != nil {
			var er ErrorResponse
			er.Code = fiber.StatusInternalServerError
			er.Data = err.Error()
			return c.Status(fiber.StatusInternalServerError).JSON(er)
		}

		var sr SucccessResponse
		sr.Code = fiber.StatusCreated
		sr.Data = "Success to created a flight!"
		return c.Status(fiber.StatusCreated).JSON(sr)
	}
}

func getAllSeats(c *fiber.Ctx) error {
	rows, err := DB.Query("SELECT id, flight_id, label, cabin, is_assigned FROM seats")
	if err != nil {
		var er ErrorResponse
		er.Code = fiber.StatusInternalServerError
		er.Data = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(er)
	}

	var seats []Seat
	for rows.Next() {
		var seat Seat
		if err := rows.Scan(&seat.ID, &seat.FlightID, &seat.Label, &seat.Cabin, &seat.IsAssigned); err != nil {
			var er ErrorResponse
			er.Code = fiber.StatusInternalServerError
			er.Data = err.Error()
			return c.Status(fiber.StatusInternalServerError).JSON(er)
		}
		seats = append(seats, seat)
	}

	var sr SucccessResponse
	sr.Code = fiber.StatusOK
	sr.Data = seats
	return c.Status(fiber.StatusOK).JSON(sr)
}

func createNewSeat(c *fiber.Ctx) error {
	sp := new(SeatPayload)
	if err := c.BodyParser(&sp); err != nil {
		var er ErrorResponse
		er.Code = fiber.StatusBadRequest
		er.Data = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(er)
	}

	sp.Cabin = strings.ToUpper(strings.TrimSpace(sp.Cabin))
	if !sp.validate() {
		var er ErrorResponse
		er.Code = fiber.StatusBadRequest
		er.Data = "Invalid JSON"
		return c.Status(fiber.StatusBadRequest).JSON(er)
	} else {
		tx, _ := DB.Begin()
		defer tx.Rollback()

		for _, lb := range sp.Labels {
			lb = strings.ToUpper(strings.TrimSpace(lb))
			if _, err := tx.Exec(`INSERT INTO seats(flight_id, label, cabin) VALUES(?,?,?)`, sp.FlightID, lb, sp.Cabin); err != nil {
				var er ErrorResponse
				er.Code = fiber.StatusInternalServerError
				er.Data = err.Error()
				return c.Status(fiber.StatusInternalServerError).JSON(er)
			}
		}

		if err := tx.Commit(); err != nil {
			var er ErrorResponse
			er.Code = fiber.StatusInternalServerError
			er.Data = err.Error()
			return c.Status(fiber.StatusInternalServerError).JSON(er)
		}

		var sr SucccessResponse
		sr.Code = fiber.StatusCreated
		sr.Data = "Success to created new seat!"
		return c.Status(fiber.StatusCreated).JSON(sr)
	}
}

func getAllVouchers(c *fiber.Ctx) error {
	rows, err := DB.Query("SELECT id, flight_id, code, cabin, redeemed, expires_at, redeemed_at FROM vouchers")
	if err != nil {
		var er ErrorResponse
		er.Code = fiber.StatusInternalServerError
		er.Data = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(er)
	}

	var vouchers []Voucher
	for rows.Next() {
		var voucher Voucher
		if err := rows.Scan(&voucher.ID, &voucher.FlightID, &voucher.Code, &voucher.Cabin, &voucher.Redeemed, &voucher.ExpiresAt, &voucher.ReadeemedAt); err != nil {
			var er ErrorResponse
			er.Code = fiber.StatusInternalServerError
			er.Data = err.Error()
			return c.Status(fiber.StatusInternalServerError).JSON(er)
		}
		vouchers = append(vouchers, voucher)
	}

	var sr SucccessResponse
	sr.Code = fiber.StatusOK
	sr.Data = vouchers
	return c.Status(fiber.StatusOK).JSON(sr)
}

func createNewVoucher(c *fiber.Ctx) error {
	vp := new(VoucherPayload)
	if err := c.BodyParser(&vp); err != nil {
		var er ErrorResponse
		er.Code = fiber.StatusBadRequest
		er.Data = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(er)
	}

	vp.Cabin = strings.ToUpper(strings.TrimSpace(vp.Cabin))
	if !vp.validate() {
		var er ErrorResponse
		er.Code = fiber.StatusBadRequest
		er.Data = "Invalid JSON"
		return c.Status(fiber.StatusBadRequest).JSON(er)
	} else {
		if vp.ExpiresAt.String() == "" {
			vp.ExpiresAt = time.Now()
		}

		if _, err := DB.Exec(`INSERT INTO vouchers(code, flight_id, cabin) VALUES(?, ?, ?)`, vp.Code, vp.FlightID, vp.Cabin); err != nil {
			var er ErrorResponse
			er.Code = fiber.StatusInternalServerError
			er.Data = err.Error()
			return c.Status(fiber.StatusInternalServerError).JSON(er)
		}

		var sr SucccessResponse
		sr.Code = fiber.StatusCreated
		sr.Data = "Success to create a new voucher!"
		return c.Status(fiber.StatusCreated).JSON(sr)
	}
}

// @TODO
// - [] Validate a voucers
// - [] Validate a seats
// - [] Add new assigments
func AddNewAssigments(c *fiber.Ctx) error {
	var sr SucccessResponse
	sr.Code = fiber.StatusCreated
	sr.Data = "Success to create a new assigments"
	return c.Status(fiber.StatusCreated).JSON(sr)
}

func main() {
	cfg := LoadConfig()

	initDB(cfg.DBPath)
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
	seats.Post("/", createNewSeat)

	// vouchers
	vouchers := v1.Group("/vouchers")
	vouchers.Get("/", getAllVouchers)
	vouchers.Post("/", createNewVoucher)

	// flights
	flights := v1.Group("/flights")
	flights.Get("/", getAllFlights)
	flights.Post("/", createNewFlight)

	// assignments
	assigments := v1.Group("/assigments")
	assigments.Post("/", AddNewAssigments)

	app.Listen(fmt.Sprintf(":%s", cfg.Port))
}
