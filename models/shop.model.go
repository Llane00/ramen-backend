package models

import (
	"github.com/google/uuid"
)

type Shop struct {
	Base
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null"`
	Owner       User      `gorm:"foreignKey:OwnerID"`
}

type Product struct {
	Base
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	Price       int64     `gorm:"type:bigint;not null"` // Price in cents
	Stock       int       `gorm:"not null"`
	ShopID      uuid.UUID `gorm:"type:uuid;not null"`
	Shop        Shop      `gorm:"foreignKey:ShopID"`
}

type Order struct {
	Base
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	User       User      `gorm:"foreignKey:UserID"`
	ShopID     uuid.UUID `gorm:"type:uuid;not null"`
	Shop       Shop      `gorm:"foreignKey:ShopID"`
	TotalPrice int64     `gorm:"type:bigint;not null"` // Total price in cents
	Status     string    `gorm:"type:varchar(50);not null"`
}

type OrderItem struct {
	Base
	OrderID      uuid.UUID `gorm:"type:uuid;not null"`
	Order        Order     `gorm:"foreignKey:OrderID"`
	ProductID    uuid.UUID `gorm:"type:uuid;not null"`
	Product      Product   `gorm:"foreignKey:ProductID"`
	ProductName  string    `gorm:"type:varchar(255);not null"` // Product name at the time of order
	ProductPrice int64     `gorm:"type:bigint;not null"`       // Product price at the time of order (in cents)
	Quantity     int       `gorm:"not null"`
	TotalPrice   int64     `gorm:"type:bigint;not null"` // Subtotal (quantity * unit price, in cents)
}

type Payment struct {
	Base
	OrderID       uuid.UUID `gorm:"type:uuid;not null"`
	Order         Order     `gorm:"foreignKey:OrderID"`
	UserID        uuid.UUID `gorm:"type:uuid;not null"`
	User          User      `gorm:"foreignKey:UserID"`
	Amount        int64     `gorm:"type:bigint;not null"` // Amount in cents
	PaymentMethod string    `gorm:"type:varchar(50);not null"`
	Status        string    `gorm:"type:varchar(50);not null"`
}
