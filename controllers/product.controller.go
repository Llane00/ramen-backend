package controllers

import (
	"net/http"

	"github.com/Llane00/ramen-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductController struct {
	DB *gorm.DB
}

func NewProductController(DB *gorm.DB) ProductController {
	return ProductController{DB}
}

// CreateProduct creates a new product
func (pc *ProductController) CreateProduct(ctx *gin.Context) {
	var input models.CreateProductInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shopID, err := uuid.Parse(ctx.Param("shopId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
		return
	}

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		ShopID:      shopID,
	}

	if err := pc.DB.Create(&product).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": product})
}

// GetProduct retrieves a product by its ID
func (pc *ProductController) GetProduct(ctx *gin.Context) {
	productID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := pc.DB.First(&product, productID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": product})
}

// UpdateProduct updates a product
func (pc *ProductController) UpdateProduct(ctx *gin.Context) {
	productID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := pc.DB.First(&product, productID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input models.UpdateProductInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pc.DB.Model(&product).Updates(input)

	ctx.JSON(http.StatusOK, gin.H{"data": product})
}

// DeleteProduct deletes a product
func (pc *ProductController) DeleteProduct(ctx *gin.Context) {
	productID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := pc.DB.Delete(&models.Product{}, productID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "Product deleted successfully"})
}

// ListProducts lists all products for a shop
func (pc *ProductController) ListProducts(ctx *gin.Context) {
	shopID, err := uuid.Parse(ctx.Param("shopId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
		return
	}

	var products []models.Product
	if err := pc.DB.Where("shop_id = ?", shopID).Find(&products).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list products"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": products})
}

// UpdateProductStock updates the stock of a product
func (pc *ProductController) UpdateProductStock(ctx *gin.Context) {
	productID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var input struct {
		Stock int `json:"stock" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := pc.DB.First(&product, productID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	product.Stock = input.Stock
	if err := pc.DB.Save(&product).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": product})
}
