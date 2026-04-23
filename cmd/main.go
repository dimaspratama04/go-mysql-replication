package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/rnd/mysql-replication/config"
	"github.com/rnd/mysql-replication/internal/handler"
	"github.com/rnd/mysql-replication/internal/middleware"
	"github.com/rnd/mysql-replication/internal/repository"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	_ = godotenv.Load()

	// Logger setup
	zerolog.TimeFieldFormat = time.RFC3339
	if os.Getenv("APP_ENV") != "production" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05",
		})
	}

	log.Info().Msg("Starting MySQL Replication RnD - Products API")

	// DB
	// primaryDBCfg := config.PrimaryDBConfig()
	// replicaDBCfg := config.ReplicaDBConfig()
	databaseConfig := config.ProxyDBConfig()

	db, err := config.NewInitDatabase(databaseConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect database")
	}

	// Layers
	productRepo := repository.NewProductRepository(db)
	productHandler := handler.NewProductHandler(productRepo)

	// Fiber
	app := fiber.New(fiber.Config{
		AppName:      "MySQL Replication RnD",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(middleware.RequestLogger())

	// Health check
	// app.Get("/health", func(c *fiber.Ctx) error {
	// 	sqlDB, err := db.DB()
	// 	if err != nil || sqlDB.PingContext(c.Context()) != nil {
	// 		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
	// 			"status":   "unhealthy",
	// 			"database": "unreachable",
	// 		})
	// 	}
	// 	return c.JSON(fiber.Map{
	// 		"status":   "healthy",
	// 		"database": "connected",
	// 		"service":  "mysql-replication-rnd",
	// 	})
	// })

	// Routes
	api := app.Group("/api")
	products := api.Group("/products")
	{
		products.Get("/", productHandler.GetAll)
		products.Get("/:id", productHandler.GetByID)
		products.Post("/", productHandler.Create)
		products.Put("/:id", productHandler.Update)
		products.Delete("/:id", productHandler.Delete)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info().Str("port", port).Msg("Server listening")
		if err := app.Listen(":" + port); err != nil {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	<-quit
	log.Info().Msg("Shutting down server...")
	_ = app.Shutdown()
	log.Info().Msg("Server stopped")
}
