package main

import (
	"github.com/Ex6linz/OMS/rbac-service/internal/database"
	"github.com/Ex6linz/OMS/rbac-service/internal/handlers"
	"github.com/Ex6linz/OMS/rbac-service/internal/models"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Database error:", err)
	}

	if err := db.AutoMigrate(&models.Role{}, &models.Permission{}); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	api := app.Group("/api/v1")

	// Role routes
	api.Get("/roles", handlers.GetRoles)
	api.Post("/roles", handlers.CreateRole)
	api.Get("/roles/:id", handlers.GetRoleByID)
	api.Put("/roles/:id", handlers.UpdateRole)
	api.Delete("/roles/:id", handlers.DeleteRole)

	// Permission routes
	api.Get("/permissions", handlers.GetPermissions)
	api.Post("/permissions", handlers.CreatePermission)
	api.Put("/permissions/:id", handlers.UpdatePermission)
	api.Delete("/permissions/:id", handlers.DeletePermission)

	// Health endpoint
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "RBAC Service running!",
		})
	})

	log.Fatal(app.Listen(":4000"))
}
