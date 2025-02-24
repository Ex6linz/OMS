package service

import (
	"github.com/Ex6linz/OMS/order-service/internal/models"
	"github.com/Ex6linz/OMS/order-service/internal/repository"
)

type OrderService interface {
	CreateOrder(order *models.Order) error
	GetOrders(limit, offset int, orderID, productName string) ([]models.Order, error)
	GetOrderByID(id uint) (*models.Order, error)
	UpdateOrder(order *models.Order) error
	DeleteOrder(id uint) error
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateOrder(order *models.Order) error {
	return s.repo.Create(order)
}

func (s *orderService) GetOrders(limit, offset int, orderID, productName string) ([]models.Order, error) {
	return s.repo.GetAll(limit, offset, orderID, productName)
}

func (s *orderService) GetOrderByID(id uint) (*models.Order, error) {
	return s.repo.GetByID(id)
}

func (s *orderService) UpdateOrder(order *models.Order) error {
	return s.repo.Update(order)
}

func (s *orderService) DeleteOrder(id uint) error {
	return s.repo.Delete(id)
}
