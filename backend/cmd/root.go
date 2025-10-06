package cmd

import (
	"backend/config"
	"backend/delivery/http"
	"backend/delivery/http/handler"
	"backend/delivery/http/middleware"
	"backend/internal/controller"
	"backend/internal/repository"
	"backend/pkg/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var server = &cobra.Command{
	Use:   "server",
	Short: "Run a web server",
	Run: func(cmd *cobra.Command, args []string) {
		// init config
		cfg := config.LoadConfig()

		// init open connection
		sqlConnection, err := db.NewSQLiteConnection(cfg.DBPath)
		if err != nil {
			log.Fatalf("Failed to init database connection: %v", err)
		}

		// execute to insert database schema
		if _, err := sqlConnection.Exec(db.SCHEMA); err != nil {
			log.Fatalf("Failed to init database schema: %v", err)
		}
		log.Info("Successfully initialized database schema!")

		// init fiber
		app := fiber.New(fiber.Config{})

		// middleware modules
		middleware.Middleware(app)

		// repository (data layer)
		flightsRepository := repository.NewFlightsRepository(sqlConnection)
		seatsRepository := repository.NewSeatRepository(sqlConnection)
		vouchersRepository := repository.NewVouchersRepository(sqlConnection)

		// controller (business layer)
		flightsController := controller.NewFlightsController(flightsRepository)
		seatsController := controller.NewSeatController(seatsRepository)
		vouchersController := controller.NewVouchersController(vouchersRepository)

		// handler (presentation layer)
		flightsHandler := handler.NewFlightsHandler(flightsController)
		seatsHandler := handler.NewSeatsHandler(seatsController)
		vouchersHandler := handler.NewVouchersHandler(vouchersController)

		// setup routes
		http.Routes(app, flightsHandler, seatsHandler, vouchersHandler)

		app.Listen(":" + cfg.Port)
	},
}

func init() {
	rootCmd.AddCommand(server)
}

func Exec() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
