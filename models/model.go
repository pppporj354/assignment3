package models

import (
	"gorm.io/gorm"
	"time"
)

type Item struct {
	gorm.Model
	ItemID      uint   `gorm:"primaryKey"`
	ItemCode    string `gorm:"size:255"`
	Description string `gorm:"size:255"`
	Quantity    int
	OrderID     uint
}

type Order struct {
	gorm.Model
	OrderID      uint   `gorm:"primaryKey"`
	CustomerName string `gorm:"size:255"`
	OrderedAt    time.Time
	Items        []Item `gorm:"foreignKey:OrderID"`
}
