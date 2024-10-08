package routes

import (
	"github.com/Llane00/ramen-backend/controllers"
	"github.com/Llane00/ramen-backend/middleware"
	"github.com/gin-gonic/gin"
)

type ProductRouteController struct {
	productController controllers.ProductController
}

func NewProductRouteController(productController controllers.ProductController) ProductRouteController {
	return ProductRouteController{productController}
}

func (pc *ProductRouteController) ProductRoute(rg *gin.RouterGroup) {
	router := rg.Group("/shops/:shopId/products")
	router.Use(middleware.DeserializeUser())

	router.POST("/", pc.productController.CreateProduct)
	router.GET("/", pc.productController.ListProducts)
	router.GET("/:productId", pc.productController.GetProduct)
	router.PUT("/:productId", pc.productController.UpdateProduct)
	router.DELETE("/:productId", pc.productController.DeleteProduct)
	router.PATCH("/:productId/stock", pc.productController.UpdateProductStock)
}
