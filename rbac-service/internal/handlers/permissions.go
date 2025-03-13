package handlers

import (
	"github.com/Ex6linz/OMS/rbac-service/internal/database"
	"github.com/Ex6linz/OMS/rbac-service/internal/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// GetPermissions - lista wszystkich permisji
func GetPermissions(c *fiber.Ctx) error {
	db, _ := database.Connect()
	var perms []models.Permission
	db.Find(&perms)
	return c.JSON(perms)
}

// CreatePermission - tworzenie permisji
func CreatePermission(c *fiber.Ctx) error {
	db, _ := database.Connect()
	perm := new(models.Permission)
	if err := c.BodyParser(&perm); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	db.Create(&perm)
	return c.Status(201).JSON(perm)
}

// UpdatePermission - edycja permisji
func UpdatePermission(c *fiber.Ctx) error {
	db, _ := database.Connect()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var perm models.Permission
	if db.First(&perm, id).Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Permission not found"})
	}

	if err := c.BodyParser(&perm); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	db.Save(&perm)
	return c.JSON(perm)
}

// DeletePermission - usunięcie permisji
func DeletePermission(c *fiber.Ctx) error {
	db, _ := database.Connect()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	db.Delete(&models.Permission{}, id)
	return c.SendStatus(204)
}
