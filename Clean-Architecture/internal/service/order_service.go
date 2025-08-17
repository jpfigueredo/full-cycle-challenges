package service

import (
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/domain"
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/repository"
)

type OrderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) GetOrders() ([]domain.Order, error) {
	return s.repo.FindAll()
}

func (s *OrderService) GetOrderByID(id uint) (*domain.Order, error) {
	return s.repo.FindByID(id)
}

func (s *OrderService) CreateOrder(order *domain.Order) error {
	return s.repo.Create(order)
}

func (s *OrderService) UpdateOrder(order *domain.Order) error {
	return s.repo.Update(order)
}

func (s *OrderService) DeleteOrder(id uint) error {
	return s.repo.Delete(id)
}
