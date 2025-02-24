package repository

import (
	"github.com/Ex6linz/OMS/order-service/internal/models"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetAll(limit, offset int, orderID, productName string) ([]models.Order, error) {
	var orders []models.Order
	query := r.db.Limit(limit).Offset(offset)
	if orderID != "" {
		query = query.Where("id = ?", orderID)
	}
	if productName != "" {
		query = query.Where("product_name ILIKE ?", "%"+productName+"%")
	}
	err := query.Find(&orders).Error
	return orders, err
}

func (r *orderRepository) GetByID(id uint) (*models.Order, error) {
	var order models.Order
	if err := r.db.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Order{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
