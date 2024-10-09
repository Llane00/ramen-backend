package controllers

import (
	"net/http"

	"github.com/Llane00/ramen-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShopController struct {
	DB *gorm.DB
}

func NewShopController(DB *gorm.DB) ShopController {
	return ShopController{DB}
}

// CreateShop creates a new shop
func (sc *ShopController) CreateShop(ctx *gin.Context) {
	var input models.CreateShopInput
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

	shop := models.Shop{
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     currentUser.ID,
	}

	if err := sc.DB.Create(&shop).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shop"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": shop})
}

// GetShop retrieves a shop by its ID
func (sc *ShopController) GetShop(ctx *gin.Context) {
	shopID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
		return
	}

	var shop models.Shop
	if err := sc.DB.First(&shop, shopID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Shop not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": shop})
}

// UpdateShop updates a shop
func (sc *ShopController) UpdateShop(ctx *gin.Context) {
	shopID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
		return
	}

	var shop models.Shop
	if err := sc.DB.First(&shop, shopID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Shop not found"})
		return
	}

	var input models.UpdateShopInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sc.DB.Model(&shop).Updates(input)

	ctx.JSON(http.StatusOK, gin.H{"data": shop})
}

// DeleteShop deletes a shop
func (sc *ShopController) DeleteShop(ctx *gin.Context) {
	shopID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
		return
	}

	if err := sc.DB.Delete(&models.Shop{}, shopID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete shop"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "Shop deleted successfully"})
}

// ListShops lists all shops
func (sc *ShopController) ListShops(ctx *gin.Context) {
	var shops []models.Shop
	if err := sc.DB.Find(&shops).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list shops"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": shops})
}

// GetShopProducts retrieves all products for a specific shop
func (sc *ShopController) GetShopProducts(ctx *gin.Context) {
	shopID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
		return
	}

	var products []models.Product
	if err := sc.DB.Where("shop_id = ?", shopID).Find(&products).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve shop products"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": products})
}

// GetShopOrders retrieves all orders for a specific shop
func (sc *ShopController) GetShopOrders(ctx *gin.Context) {
	shopID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
		return
	}

	var orders []models.Order
	if err := sc.DB.Where("shop_id = ?", shopID).Find(&orders).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve shop orders"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": orders})
}
