package models

type Order struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	CustomerID  uint   `json:"customer_id"`
}