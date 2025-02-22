package main

import (
	"github.com/Ex6linz/OMS/order-service/internal/handlers"
	"github.com/Ex6linz/OMS/order-service/internal/logger"
	"github.com/Ex6linz/OMS/order-service/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

func main() {


	logger.Init()
	defer logger.Log.Sync()


	if err := godotenv.Load(); err != nil {
		logger.Log.Warn("Brak pliku .env", zap.Error(err))
	}


	app := fiber.New()

	// Middleware logujące
	app.Use(func(c *fiber.Ctx) error {
		logger.Log.Info("Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)
		return c.Next()
	})

	// Routing
	api := app.Group("/api/v1")
	orders := api.Group("/orders")

	// Middleware JWT
	orders.Use(middleware.NewAuthMiddleware(os.Getenv("JWT_SECRET")))

	// Endpointy - usuniete przecinki i dodane średniki
	orders.Post("/", handlers.CreateOrder)
	orders.Get("/", handlers.GetAllOrders)
	orders.Get("/:id", handlers.GetOrderByID)
	orders.Put("/:id", handlers.UpdateOrder)
	orders.Delete("/:id", handlers.DeleteOrder)

	// Start serwera
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	logger.Log.Info("Starting server", zap.String("port", port))
	if err := app.Listen(":" + port); err != nil {
		logger.Log.Fatal("Server error", zap.Error(err))
	}
}
