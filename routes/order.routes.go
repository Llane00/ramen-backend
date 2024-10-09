package routes

import (
	"github.com/Llane00/ramen-backend/controllers"
	"github.com/Llane00/ramen-backend/middleware"
	"github.com/gin-gonic/gin"
)

type OrderRouteController struct {
	orderController controllers.OrderController
}

func NewOrderRouteController(orderController controllers.OrderController) OrderRouteController {
	return OrderRouteController{orderController}
}

func (oc *OrderRouteController) OrderRoute(rg *gin.RouterGroup) {
	router := rg.Group("/shops/:shopId/orders")
	router.Use(middleware.DeserializeUser())

	router.POST("/", oc.orderController.CreateOrder)
	router.GET("/", oc.orderController.ListOrders)
	router.GET("/:id", oc.orderController.GetOrder)
	router.PATCH("/:id/status", oc.orderController.UpdateOrderStatus)
	router.GET("/:id/payments", oc.orderController.GetOrderPayments)
}
