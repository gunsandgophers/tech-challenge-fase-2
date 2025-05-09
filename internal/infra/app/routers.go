package app

import (
	"tech-challenge-fase-1/internal/infra/controllers"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Registra as rotas dos controllers
func registerRouters(app *APIApp) {
	helloController := controllers.NewHelloController()
	// customerController := controllers.NewCustomerController(app.customerRepository)
	productController := controllers.NewProductController(app.productRepository)
	orderController := controllers.NewOrderController(
		app.orderRepository,
		app.customerService,
		app.productRepository,
		app.mercadoPagoGateway,
		app.eventManager,
		app.orderDisplayListQuery,
	)

	baseUrl := "/api/v1"
	app.httpServer.SetBasePath(baseUrl)
	app.httpServer.GET("/", helloController.Index)

	//customer
	// app.httpServer.POST("/customer/", customerController.RegisterCustomer)
	// app.httpServer.GET("/customer/:cpf/", customerController.GetCustomer)

	//products
	app.httpServer.POST("/product", productController.CreateProduct)
	app.httpServer.PUT("/product/:id", productController.UpdateProduct)
	app.httpServer.DELETE("/product/:id", productController.DeleteProduct)
	app.httpServer.GET("/product/:category", productController.ListProductsByCategory)

	//orders
	app.httpServer.POST("/order/checkout", orderController.Checkout)
	app.httpServer.GET(
		"/order/:order_id/payment-status",
		orderController.GetPaymentStatus,
	)
	app.httpServer.POST("/order/payment", orderController.Payment)
	app.httpServer.GET("/order/display", orderController.OrderDisplayList)
	app.httpServer.PUT(
		"/order/:order_id/preparation-status",
		orderController.OrderPreparationStatusUpdate,
	)

	app.httpServer.SetSwagger("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
