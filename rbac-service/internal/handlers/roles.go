package handlers

import (
	"strconv"

	"github.com/Ex6linz/OMS/rbac-service/internal/database"
	"github.com/Ex6linz/OMS/rbac-service/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetRoles(c *fiber.Ctx) error {
	db, _ := database.Connect()
	var roles []models.Role
	db.Preload("Permissions").Find(&roles)
	return c.JSON(roles)
}

func CreateRole(c *fiber.Ctx) error {
	db, _ := database.Connect()
	role := new(models.Role)

	if err := c.BodyParser(role); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	db.Create(&role)
	return c.Status(201).JSON(role)
}

func GetRoleByID(c *fiber.Ctx) error {
	db, _ := database.Connect()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var role models.Role
	if result := db.Preload("Permissions").First(&role, id); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Role not found"})
	}

	return c.JSON(role)
}

func UpdateRole(c *fiber.Ctx) error {
	db, _ := database.Connect()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var role models.Role
	if db.First(&role, id).Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Role not found"})
	}

	if err := c.BodyParser(&role); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	db.Save(&role)
	return c.JSON(role)
}

func DeleteRole(c *fiber.Ctx) error {
	db, _ := database.Connect()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	db.Delete(&models.Role{}, id)
	return c.SendStatus(204)
}
