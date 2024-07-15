package orders

import (
	"tech-challenge-fase-1/internal/core/dtos"
	"tech-challenge-fase-1/internal/core/entities"
	"tech-challenge-fase-1/internal/core/events"
	"tech-challenge-fase-1/internal/core/repositories"
	"tech-challenge-fase-1/internal/core/services"
)

type CheckoutOrderUseCase struct {
	orderRepository     repositories.OrderRepositoryInterface
	customerRepository repositories.CustomerRepositoryInterface
	productRepository repositories.ProductRepositoryInterface
	paymentGateway      services.PaymentGatewayInterface
	commandEventManager events.ManagerEvent
}

func NewCheckoutOrderUseCase(
	orderRepository repositories.OrderRepositoryInterface,
	customerRepository repositories.CustomerRepositoryInterface,
	productRepository repositories.ProductRepositoryInterface,
	paymentGateway services.PaymentGatewayInterface,
	commandEventManager events.ManagerEvent,
) *CheckoutOrderUseCase {
	return &CheckoutOrderUseCase{
		orderRepository:     orderRepository,
		customerRepository:     customerRepository,
		productRepository:     productRepository,
		paymentGateway:      paymentGateway,
		commandEventManager: commandEventManager,
	}
}

func (c *CheckoutOrderUseCase) validateCustomerId(customerId *string) error {
	if customerId == nil {
		return nil
	}
	_, err := c.customerRepository.GetCustomerByID(*customerId)
	if err != nil {
		return err
	}
	return nil
}

func (c *CheckoutOrderUseCase) fetchProducts(productsIds []string) ([]*entities.Product, error) {
	products := []*entities.Product{}
	for _, productId := range productsIds {
		product, err := c.productRepository.FindProductByID(productId)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (c *CheckoutOrderUseCase) Execute(
	customerId *string,
	productsIds []string,
) (*dtos.CheckoutDTO, error) {
	if err := c.validateCustomerId(customerId); err != nil {
		return nil, err
	}
	products, err := c.fetchProducts(productsIds)
	if err != nil {
		return nil, err
	}
	order := entities.CreateOpenOrder(customerId)
	for _, product := range products {
		order.AddItem(product, 1)
	}

	// TODO: refactor. Stay here just for now
	order.SetStatus(entities.AWAITING_PAYMENT)
	checkout, err := c.paymentGateway.Execute(
		dtos.NewOrderDTOFromEntity(order),
		dtos.PIX,
	)
	if err != nil {
		return nil, err
	}

	c.orderRepository.Insert(order)

	// TODO: refactor/remove??
	// err = c.orderRepository.Update(order)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// go func() {
	// 	c.commandEventManager.Add("paid_out", func() {
	// 		order.SetStatus(entities.PAID)
	// 		c.orderRepository.Update(order)
	// 	})
	// }()

	return checkout, nil
}
