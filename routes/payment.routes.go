package routes

import (
	"github.com/Llane00/ramen-backend/controllers"
	"github.com/Llane00/ramen-backend/middleware"
	"github.com/gin-gonic/gin"
)

type PaymentRouteController struct {
	paymentController controllers.PaymentController
}

func NewPaymentRouteController(paymentController controllers.PaymentController) PaymentRouteController {
	return PaymentRouteController{paymentController}
}

func (pc *PaymentRouteController) PaymentRoute(rg *gin.RouterGroup) {
	router := rg.Group("/orders/:orderId/payments")
	router.Use(middleware.DeserializeUser())

	router.POST("/", pc.paymentController.CreatePayment)
	router.GET("/", pc.paymentController.ListPayments)
	router.GET("/:id", pc.paymentController.GetPayment)
	router.PATCH("/:id/status", pc.paymentController.UpdatePaymentStatus)
}
