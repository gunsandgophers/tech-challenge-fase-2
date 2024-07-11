package orders

import (
	"tech-challenge-fase-1/internal/core/repositories"
	"tech-challenge-fase-1/internal/core/dtos"
)

type AddOrderItemUseCaseRequest struct {
	ProductID int
	Quantity  int
	OrderID   string
}

type AddOrderItemUseCase struct {
	orderRepository   repositories.OrderRepositoryInterface
	productRepository repositories.ProductRepositoryInterface
}

func NewAddOrderItemUseCase(
	orderRepository repositories.OrderRepositoryInterface,
	productRepository repositories.ProductRepositoryInterface,
) *AddOrderItemUseCase {
	return &AddOrderItemUseCase{
		orderRepository:   orderRepository,
		productRepository: productRepository,
	}
}

func (co *AddOrderItemUseCase) Execute(request *AddOrderItemUseCaseRequest) (*dtos.OrderDTO, error) {
	_, err := co.productRepository.FindByID(request.ProductID)
	if err != nil {
		return nil, err
	}
	order, err := co.orderRepository.FindByID(&request.OrderID)
	if err != nil {
		return nil, err
	}
	return dtos.NewOrderDTOFromEntity(order), nil
	// order.AddItem(product, &request.Quantity)
	// return dtos.NewAddOrderItemDTOFromEntity(order), co.orderRepository.AddOrderItem(order)
}
