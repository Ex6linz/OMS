package main

import (
	"github.com/Ex6linz/OMS/order-service/internal/database"
	"github.com/Ex6linz/OMS/order-service/internal/handlers"
	"github.com/Ex6linz/OMS/order-service/internal/logger"
	"github.com/Ex6linz/OMS/order-service/internal/middleware"
	"github.com/Ex6linz/OMS/order-service/internal/models"
	"github.com/Ex6linz/OMS/order-service/internal/rbac"
	_ "github.com/Ex6linz/OMS/order-service/internal/rbac"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/gofiber/fiber/v2/middleware/cors"
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

	// Uruchamiamy subskrypcję z Kafki
	// Parametry (adresy brokerów, topic, groupID) możesz pobrać z ENV
	kafkaBrokers := []string{"localhost:9092"}
	topic := "rbac_updates"
	groupID := "order-service-group"
	rbac.StartRBACConsumer(kafkaBrokers, topic, groupID)

	// Połączenie z bazą danych i migracja modeli
	db, err := database.Connect()
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	// Migracja modelu Order (utworzy tabelę, jeśli nie istnieje)
	if err := db.AutoMigrate(&models.Order{}); err != nil {
		logger.Log.Fatal("Failed to migrate database", zap.Error(err))
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
	app.Use(cors.New())
	// Routing
	api := app.Group("/api/v1")
	orders := api.Group("/orders")
	orders.Delete("/:id", handlers.DeleteOrder)
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
