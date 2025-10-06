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

// TestApp holds the test application and database
type TestApp struct {
	App *fiber.App
	DB  *sql.DB
}

// setupTestApp creates a new test application with in-memory database
func setupTestApp(t *testing.T) *TestApp {
	// Create in-memory SQLite database
	database, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Initialize schema
	if _, err := database.Exec(db.SCHEMA); err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	// Create repositories
	flightsRepo := repository.NewFlightsRepository(database)
	seatsRepo := repository.NewSeatRepository(database)
	vouchersRepo := repository.NewVouchersRepository(database)

	// Create controllers
	flightsController := controller.NewFlightsController(flightsRepo)
	seatsController := controller.NewSeatController(seatsRepo)
	vouchersController := controller.NewVouchersController(vouchersRepo)

	// Create handlers
	flightsHandler := handler.NewFlightsHandler(flightsController)
	seatsHandler := handler.NewSeatsHandler(seatsController)
	vouchersHandler := handler.NewVouchersHandler(vouchersController)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Setup routes
	http.Routes(app, flightsHandler, seatsHandler, vouchersHandler)

	return &TestApp{
		App: app,
		DB:  database,
	}
}

// cleanup closes the database connection
func (ta *TestApp) cleanup() {
	if ta.DB != nil {
		ta.DB.Close()
	}
}

// makeRequest is a helper function to make HTTP requests to the test app
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

	// Convert to ResponseRecorder for easier testing
	recorder := httptest.NewRecorder()
	recorder.WriteHeader(resp.StatusCode)
	io.Copy(recorder, resp.Body)
	resp.Body.Close()

	return recorder, nil
}

// parseResponse parses JSON response into the provided struct
func parseResponse(t *testing.T, resp *httptest.ResponseRecorder, v any) {
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
}
