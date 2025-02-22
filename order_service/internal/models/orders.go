package models

type Order struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ProductName string `json:"product_name" validate:"required,min=3"`
	Quantity    int    `json:"quantity" validate:"required,gt=0"`
	CustomerID  uint   `json:"customer_id" validate:"required"`
}