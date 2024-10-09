package controllers

import (
	"net/http"

	"github.com/Llane00/ramen-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentController struct {
	DB *gorm.DB
}

func NewPaymentController(DB *gorm.DB) PaymentController {
	return PaymentController{DB}
}

// CreatePayment creates a new payment
func (pc *PaymentController) CreatePayment(ctx *gin.Context) {
	var input models.CreatePaymentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID, err := uuid.Parse(ctx.Param("orderId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	if err := pc.DB.First(&order, orderID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	payment := models.Payment{
		OrderID:       orderID,
		Amount:        input.Amount,
		PaymentMethod: input.PaymentMethod,
		Status:        models.PaymentStatusPending,
	}

	if err := pc.DB.Create(&payment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": payment})
}

// GetPayment retrieves a payment by its ID
func (pc *PaymentController) GetPayment(ctx *gin.Context) {
	paymentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	var payment models.Payment
	if err := pc.DB.First(&payment, paymentID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": payment})
}

// UpdatePaymentStatus updates the status of a payment
func (pc *PaymentController) UpdatePaymentStatus(ctx *gin.Context) {
	paymentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	var input struct {
		Status models.PaymentStatus `json:"status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var payment models.Payment
	if err := pc.DB.First(&payment, paymentID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	payment.Status = input.Status
	if err := pc.DB.Save(&payment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": payment})
}

// ListPayments lists all payments for an order
func (pc *PaymentController) ListPayments(ctx *gin.Context) {
	orderID, err := uuid.Parse(ctx.Param("orderId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var payments []models.Payment
	if err := pc.DB.Where("order_id = ?", orderID).Find(&payments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list payments"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": payments})
}
