package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ex6linz/OMS/order-service/internal/models"
	_ "github.com/Ex6linz/OMS/order-service/internal/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Nie udało się otworzyć bazy: %v", err)
	}
	if err := db.AutoMigrate(&models.Order{}); err != nil {
		t.Fatalf("Nie udało się zmigrować bazy: %v", err)
	}
	return db
}

func TestGetAllOrdersNoPermission(t *testing.T) {
	app := fiber.New()

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("permissions", []string{})
		return c.Next()
	})

	// Rejestracja testowego handlera
	app.Get("/orders", GetAllOrders)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Błąd testu: %v", err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Oczekiwano status 403, otrzymano %d", resp.StatusCode)
	}
}

func TestCreateOrderWithPermission(t *testing.T) {
	db = setupTestDB(t)

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("permissions", []string{"create_order"})
		return c.Next()
	})
	app.Post("/orders", CreateOrder)

	payload := `{"product_name": "Testowy produkt", "quantity": 5, "customer_id": 1}`
	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Błąd testu: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Oczekiwano status 201, otrzymano %d", resp.StatusCode)
	}
}
