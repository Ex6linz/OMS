package database

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)




func Connect() (*g0rm.DB, error) {
	dsn := "host=postgres user=postgres password=secret dbname=orders port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    return db, err
}