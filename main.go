package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Llane00/ramen-backend/controllers"
	"github.com/Llane00/ramen-backend/initializers"
	"github.com/Llane00/ramen-backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	PostController      controllers.PostController
	PostRouteController routes.PostRouteController

	ShopController      controllers.ShopController
	ShopRouteController routes.ShopRouteController

	ProductController      controllers.ProductController
	ProductRouteController routes.ProductRouteController

	OrderController      controllers.OrderController
	OrderRouteController routes.OrderRouteController

	PaymentController      controllers.PaymentController
	PaymentRouteController routes.PaymentRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development" // default env
	}
	initializers.ConnectDB(&config, env)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(initializers.DB)
	UserRouteController = routes.NewRouteUserController(UserController)

	PostController = controllers.NewPostController(initializers.DB)
	PostRouteController = routes.NewRoutePostController(PostController)

	ShopController = controllers.NewShopController(initializers.DB)
	ShopRouteController = routes.NewShopRouteController(ShopController)

	ProductController = controllers.NewProductController(initializers.DB)
	ProductRouteController = routes.NewProductRouteController(ProductController)

	OrderController = controllers.NewOrderController(initializers.DB)
	OrderRouteController = routes.NewOrderRouteController(OrderController)

	PaymentController = controllers.NewPaymentController(initializers.DB)
	PaymentRouteController = routes.NewPaymentRouteController(PaymentController)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and Postgres"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	PostRouteController.PostRoute(router)
	ShopRouteController.ShopRoute(router)
	ProductRouteController.ProductRoute(router)
	OrderRouteController.OrderRoute(router)
	PaymentRouteController.PaymentRoute(router)
	log.Fatal(server.Run(":" + config.ServerPort))
}
