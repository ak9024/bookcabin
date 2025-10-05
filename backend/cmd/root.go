package cmd

import (
	"backend/config"
	"backend/delivery/http"
	"backend/delivery/http/handler"
	"backend/delivery/http/middleware"
	"backend/internal/controller"
	"backend/internal/repository"
	"backend/pkg/db"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
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
			log.Println("Error to init database connection!")
		}

		// execute to insert database schema
		if _, err := sqlConnection.Exec(db.SCHEMA); err != nil {
			log.Println("Error to init database schema!")
		} else {
			log.Println("Success to insert database schema!")
		}

		// repository (data layer)
		flightsRepository := repository.NewFlightsRepository(sqlConnection)

		// controller (business layer)
		flightsController := controller.NewFlightsController(flightsRepository)

		// handler (presentation layer)
		flightsHandler := handler.NewFlightsHandler(flightsController)
		seatsHandler := handler.NewSeatsHandler()
		vouchersHandler := handler.NewVouchersHandler()
		assignmentsHandler := handler.NewAssignmentsHandler()

		// init fiber
		app := fiber.New(fiber.Config{})

		// middleware modules
		middleware.Middleware(app)

		// setup routes
		http.Routes(app, flightsHandler, seatsHandler, vouchersHandler, assignmentsHandler)

		app.Listen(":" + cfg.Port)
	},
}

func init() {
	rootCmd.AddCommand(server)
}

func Exec() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
