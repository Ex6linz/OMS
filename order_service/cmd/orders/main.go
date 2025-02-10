package main

import (
	"github.com/Ex6linz/OMS/order-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    // Middleware (opcjonalnie)
    app.Use(func(c *fiber.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.Next()
    })

    // Routing
    app.Post("/orders", handlers.CreateOrder)
    app.Get("/orders", handlers.GetAllOrders)
    app.Get("/orders/:id", handlers.GetOrderByID)
    app.Put("/orders/:id", handlers.UpdateOrder)
    app.Delete("/orders/:id", handlers.DeleteOrder)

    app.Listen(":3000")
}