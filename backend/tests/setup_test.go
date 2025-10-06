package tests

import (
	"backend/delivery/http"
	"backend/delivery/http/handler"
	"backend/internal/controller"
	"backend/internal/repository"
	"backend/pkg/db"
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
)

type TestApp struct {
	App *fiber.App
	DB  *sql.DB
}

func setupTestApp(t *testing.T) *TestApp {
	database, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	if _, err := database.Exec(db.SCHEMA); err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	flightsRepo := repository.NewFlightsRepository(database)
	seatsRepo := repository.NewSeatRepository(database)
	vouchersRepo := repository.NewVouchersRepository(database)

	flightsController := controller.NewFlightsController(flightsRepo)
	seatsController := controller.NewSeatController(seatsRepo)
	vouchersController := controller.NewVouchersController(vouchersRepo)

	flightsHandler := handler.NewFlightsHandler(flightsController)
	seatsHandler := handler.NewSeatsHandler(seatsController)
	vouchersHandler := handler.NewVouchersHandler(vouchersController)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	http.Routes(app, flightsHandler, seatsHandler, vouchersHandler)

	return &TestApp{
		App: app,
		DB:  database,
	}
}

func (ta *TestApp) cleanup() {
	if ta.DB != nil {
		ta.DB.Close()
	}
}

func (ta *TestApp) makeRequest(method, path string, body any) (*httptest.ResponseRecorder, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	resp, err := ta.App.Test(req, -1) // -1 disables timeout
	if err != nil {
		return nil, err
	}

	recorder := httptest.NewRecorder()
	recorder.WriteHeader(resp.StatusCode)
	io.Copy(recorder, resp.Body)
	resp.Body.Close()

	return recorder, nil
}

func parseResponse(t *testing.T, resp *httptest.ResponseRecorder, v any) {
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
}
