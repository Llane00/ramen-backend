package controllers

import (
	"net/http"

	"github.com/Llane00/ramen-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderController struct {
	DB *gorm.DB
}

func NewOrderController(DB *gorm.DB) OrderController {
	return OrderController{DB}
}

// CreateOrder creates a new order
func (oc *OrderController) CreateOrder(ctx *gin.Context) {
	var input models.CreateOrderInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	currentUser := user.(models.User)

	shopId, err := uuid.Parse(ctx.Param("shopId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
		return
	}

	order := models.Order{
		UserID:     currentUser.ID,
		ShopID:     shopId,
		TotalPrice: input.TotalPrice,
		Status:     models.OrderStatusPending,
	}

	if err := oc.DB.Create(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": order})
}

// GetOrder retrieves an order by its ID
func (oc *OrderController) GetOrder(ctx *gin.Context) {
	orderId, err := uuid.Parse(ctx.Param("orderId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	if err := oc.DB.Preload("Items").First(&order, orderId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": order})
}

// UpdateOrderStatus updates the status of an order
func (oc *OrderController) UpdateOrderStatus(ctx *gin.Context) {
	orderId, err := uuid.Parse(ctx.Param("orderId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var input models.UpdateOrderStatusInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	if err := oc.DB.First(&order, orderId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.Status = input.Status
	if err := oc.DB.Save(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": order})
}

// ListOrders lists all orders for a shop
func (oc *OrderController) ListOrders(ctx *gin.Context) {
	shopId, err := uuid.Parse(ctx.Param("shopId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
		return
	}

	var orders []models.Order
	if err := oc.DB.Where("shop_id = ?", shopId).Find(&orders).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list orders"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": orders})
}

// GetOrderPayments retrieves all payments for a specific order
func (oc *OrderController) GetOrderPayments(ctx *gin.Context) {
	orderId, err := uuid.Parse(ctx.Param("orderId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var payments []models.Payment
	if err := oc.DB.Where("order_id = ?", orderId).Find(&payments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order payments"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": payments})
}
