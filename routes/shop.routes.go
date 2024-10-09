package routes

import (
	"github.com/Llane00/ramen-backend/controllers"
	"github.com/Llane00/ramen-backend/middleware"
	"github.com/gin-gonic/gin"
)

type ShopRouteController struct {
	shopController controllers.ShopController
}

func NewShopRouteController(shopController controllers.ShopController) ShopRouteController {
	return ShopRouteController{shopController}
}

func (sc *ShopRouteController) ShopRoute(rg *gin.RouterGroup) {
	router := rg.Group("/shops")
	router.Use(middleware.DeserializeUser())

	router.POST("/", sc.shopController.CreateShop)
	router.GET("/", sc.shopController.ListShops)
	router.GET("/:id", sc.shopController.GetShop)
	router.PUT("/:id", sc.shopController.UpdateShop)
	router.DELETE("/:id", sc.shopController.DeleteShop)
	router.GET("/:id/products", sc.shopController.GetShopProducts)
	router.GET("/:id/orders", sc.shopController.GetShopOrders)
}
