package repository

import "github.com/Ex6linz/OMS/order-service/internal/models"

type OrderRepository interface {
	Create(order *models.Order) error

	GetAll(limit, offset int, orderID, productName string) ([]models.Order, error)
	GetByID(id uint) (*models.Order, error)
	Update(order *models.Order) error
	Delete(id uint) error
}
