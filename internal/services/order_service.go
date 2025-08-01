package services

import (
	"go_ecommerce/internal/models"
	"go_ecommerce/internal/repositories"
)

type OrderService struct {
	repo *repositories.OrderRepository
}

func NewOrderService() *OrderService {
	return &OrderService{
		repo: repositories.NewOrderRepository(),
	}
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.repo.GetAllOrders()
}
func (s *OrderService) GetOrderHistory() ([]models.Order, error) {
	return s.repo.GetAllOrders()
}

