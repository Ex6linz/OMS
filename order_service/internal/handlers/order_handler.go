package handlers

import (
	"strconv"

	"github.com/Ex6linz/OMS/order-service/internal/utils"
	"github.com/go-playground/validator/v10"

	"github.com/Ex6linz/OMS/order-service/internal/database"
	"github.com/Ex6linz/OMS/order-service/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = database.Connect()
	if err != nil {
		panic("Failed to connect to database")
	}
}

func CreateOrder(c *fiber.Ctx) error {

	if !utils.HasPermission(c, "create_order") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden: insufficient permissions to create order",
		})
	}

	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := utils.Validate.Struct(order); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": errors,
		})
	}

	if result := db.Create(&order); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create order",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

func GetAllOrders(c *fiber.Ctx) error {

	if !utils.HasPermission(c, "read_orders") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden: insufficient permissions to read orders",
		})
	}

	var orders []models.Order
	// Pobieranie parametrów paginacji
	limitQuery := c.Query("limit", "10")
	pageQuery := c.Query("page", "1")
	limit, err := strconv.Atoi(limitQuery)
	if err != nil || limit <= 0 {
		limit = 10
	}
	page, err := strconv.Atoi(pageQuery)
	if err != nil || page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := db.Limit(limit).Offset(offset)

	productFilter := c.Query("product")
	if productFilter != "" {
		query = query.Where("product_name ILIKE ?", "%"+productFilter+"%")
	}

	result := db.Find(&orders)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}

	return c.Status(fiber.StatusOK).JSON(orders)
}

func GetOrderByID(c *fiber.Ctx) error {

	if !utils.HasPermission(c, "read_orders") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden: insufficient permissions to read orders",
		})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	var order models.Order
	result := db.First(&order, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(order)
}

func UpdateOrder(c *fiber.Ctx) error {

	if !utils.HasPermission(c, "update_order") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden: insufficient permissions to update orders",
		})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	var existingOrder models.Order
	if err := db.First(&existingOrder, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	var updatedOrder models.Order
	if err := c.BodyParser(&updatedOrder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	existingOrder.ProductName = updatedOrder.ProductName
	existingOrder.Quantity = updatedOrder.Quantity
	existingOrder.CustomerID = updatedOrder.CustomerID

	if err := db.Save(&existingOrder).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	return c.Status(fiber.StatusOK).JSON(existingOrder)
}

func DeleteOrder(c *fiber.Ctx) error {
	if !utils.HasPermission(c, "delete_order") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden: insufficient permissions",
		})
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	result := db.Delete(&models.Order{}, id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete order",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
