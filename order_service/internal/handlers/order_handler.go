package handlers

import (
	"strconv"

	"github.com/Ex6linz/OMS/order-service/internal/database"
	"github.com/Ex6linz/OMS/order-service/internal/models"
	"github.com/gofiber/fiber/v2"
)

func CreateOrder(c *fiber.Ctx) error {
    db, err := database.Connect()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Database error"})
    }

    var order models.Order
    if err := c.BodyParser(&order); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    result := db.Create(&order)
    if result.Error != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to create order"})
    }

    return c.JSON(order)
}

func GetOrder(c *fiber.Ctx) error {
    id := c.Params("id")
    var order models.Order
    result := db.First(&order, id)
    if result.Error != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
    }
    return c.JSON(order)
}

func GetAllOrders(c *fiber.Ctx) error {
    db, err := database.Connect()
    if err != nil {
        return c.Status(iber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
    }

    var orders []models.Order
    result := db.Find(&orders)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch orders"})
    }

    return c.JSON(orders)
}

// GetOrderByID - Pobierz pojedyncze zamówienie
func GetOrderByID(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
    }

    db, err := database.Connect()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
    }

    var order models.Order
    result := db.First(&order, id)
    if result.Error != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
    }

    return c.JSON(order)
}

// UpdateOrder - Aktualizuj zamówienie
func UpdateOrder(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
    }

    db, err := database.Connect()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
    }

    var existingOrder models.Order
    if err := db.First(&existingOrder, id).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
    }

    var updatedOrder models.Order
    if err := c.BodyParser(&updatedOrder); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
    }

    // Aktualizuj tylko dozwolone pola
    existingOrder.ProductName = updatedOrder.ProductName
    existingOrder.Quantity = updatedOrder.Quantity
    existingOrder.CustomerID = updatedOrder.CustomerID

    db.Save(&existingOrder)
    return c.JSON(existingOrder)
}

// DeleteOrder - Usuń zamówienie
func DeleteOrder(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
    }

    db, err := database.Connect()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
    }

    var order models.Order
    result := db.Delete(&order, id)
    if result.RowsAffected == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
    }

    return c.SendStatus(fiber.StatusNoContent)
}