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
	Products    []Product `gorm:"foreignKey:ShopID"`
	Orders      []Order   `gorm:"foreignKey:ShopID"`
}

type CreateShopInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateShopInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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

type CreateProductInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       int64  `json:"price" binding:"required"`
	Stock       int    `json:"stock" binding:"required"`
}

type UpdateProductInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Stock       int    `json:"stock"`
}

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipping  OrderStatus = "shipping"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	Base
	UserID     uuid.UUID   `gorm:"type:uuid;not null"`
	User       User        `gorm:"foreignKey:UserID"`
	ShopID     uuid.UUID   `gorm:"type:uuid;not null"`
	Shop       Shop        `gorm:"foreignKey:ShopID"`
	TotalPrice int64       `gorm:"type:bigint;not null"` // Total price in cents
	Status     OrderStatus `gorm:"type:varchar(50);not null"`
	Items      []OrderItem `gorm:"foreignKey:OrderID"`
	Payment    Payment     `gorm:"foreignKey:OrderID"`
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

type CreateOrderInput struct {
	TotalPrice int64                  `json:"total_price" binding:"required"`
	Items      []CreateOrderItemInput `json:"items" binding:"required,dive"`
}

type CreateOrderItemInput struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required"`
}

type UpdateOrderStatusInput struct {
	Status OrderStatus `json:"status" binding:"required"`
}

type Payment struct {
	Base
	OrderID       uuid.UUID     `gorm:"type:uuid;not null"`
	Order         *Order        `gorm:"foreignKey:OrderID"`
	Amount        int64         `gorm:"type:bigint;not null"` // Amount in cents
	PaymentMethod string        `gorm:"type:varchar(50);not null"`
	Status        PaymentStatus `gorm:"type:varchar(50);not null"`
}

type CreatePaymentInput struct {
	Amount        int64  `json:"amount" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"`
}

type UpdatePaymentStatusInput struct {
	Status string `json:"status" binding:"required"`
}
