package entities

import "time"

type Product struct {
	Id          uint `gorm:"primaryKey"`
	Name        string
	Category    Category
	Stock       int64
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
