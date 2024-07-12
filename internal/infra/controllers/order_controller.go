package controllers

import (
	"net/http"
	httpserver "tech-challenge-fase-1/internal/infra/http"
	"tech-challenge-fase-1/internal/infra/controllers/request"
	"tech-challenge-fase-1/internal/core/events"
	"tech-challenge-fase-1/internal/core/repositories"
	"tech-challenge-fase-1/internal/core/services"
	"tech-challenge-fase-1/internal/core/use_cases/orders"
)

type OrderController struct {
	orderRepository    repositories.OrderRepositoryInterface
	customerRepository repositories.CustomerRepositoryInterface
	productRepository  repositories.ProductRepositoryInterface

	paymentGateway services.PaymentGatewayInterface

	commandEventManager events.ManagerEvent
}

func NewOrderController(
	orderRepository repositories.OrderRepositoryInterface,
	customerRepository repositories.CustomerRepositoryInterface,
	productRepository repositories.ProductRepositoryInterface,
	paymentGateway services.PaymentGatewayInterface,
	commandEventManager events.ManagerEvent,
) *OrderController {
	return &OrderController{
		orderRepository:     orderRepository,
		customerRepository:  customerRepository,
		productRepository:   productRepository,
		paymentGateway:      paymentGateway,
		commandEventManager: commandEventManager,
	}
}

// OpenOrder godoc
// @Summary      Open an order
// @Description  initiate the order process
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order   body      request.OpenOrderRequest  true  "Open Order"
// @Success      200 {object} dtos.OrderDTO
// @Failure      400 {string} string "when invalid params"
// @Failure      406 {string} string "when invalid status"
// @Router       /order/open/ [post]
func (cc *OrderController) OpenOrder(c httpserver.HTTPContext) {

	request := request.OpenOrderRequest{}
	c.BindJSON(&request)

	if err := request.Validate(); err != nil {
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	openOrderUseCase := orders.NewOpenOrderUseCase(
		cc.orderRepository,
		cc.customerRepository,
	)

	order, err := openOrderUseCase.Execute(request.CustomerID)

	if err != nil {
		sendError(c, http.StatusNotAcceptable, err.Error())
		return
	}

	sendSuccess(c, http.StatusCreated, "open-order", order)
}

// AddOrderItem godoc
// @Summary      Add an order item
// @Description  insert an item to a given order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order_id   path      string true  "Open Order"
// @Param        product_id   body      int  true  "Open Order"
// @Param        quantity   body      int  true  "Open Order"
// @Success      200 {object} dtos.OrderDTO
// @Failure      400 {string} string "when invalid params"
// @Failure      406 {string} string "when invalid status"
// @Router       /order/{order_id}/add/item [post]
func (cc *OrderController) AddOrderItem(c httpserver.HTTPContext) {

	request := request.AddOrderItemRequest{}
	c.BindJSON(&request)

	request.OrderID = c.Param("order_id")

	if err := request.Validate(); err != nil {
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	addItemOrderUseCase := orders.NewAddOrderItemUseCase(
		cc.orderRepository,
		cc.productRepository,
	)

	order, err := addItemOrderUseCase.Execute(&orders.AddOrderItemUseCaseRequest{
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
		OrderID:   request.OrderID,
	})

	if err != nil {
		sendError(c, http.StatusNotAcceptable, err.Error())
		return
	}

	sendSuccess(c, http.StatusCreated, "add-order-item", order)
}

// Checkout godoc
// @Summary      Do a order checkout
// @Description  do a checkout on a given order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order_id   path      string  true  "Open Order"
// @Success      200 {object} dtos.CheckoutDTO
// @Failure      400 {string} string "when invalid status"
// @Router       /order/{order_id}/add/item [post]
func (cc *OrderController) Checkout(c httpserver.HTTPContext) {

	orderID := c.Param("order_id")

	checkoutUseCase := orders.NewCheckoutOrderUseCase(
		cc.orderRepository,
		cc.paymentGateway,
		cc.commandEventManager)

	checkout, err := checkoutUseCase.Execute(orderID)
	if err != nil {
		sendError(c, http.StatusNotAcceptable, err.Error())
		return
	}

	sendSuccess(c, http.StatusCreated, "checkout-order", checkout)

}
