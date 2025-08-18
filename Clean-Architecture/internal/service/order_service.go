package service

import (
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/domain"
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/repository"
)

type OrderService interface {
	CreateOrder(order domain.Order) (domain.Order, error)
	UpdateOrder(order *domain.Order) error
	DeleteOrder(id uint) error
	GetOrders() ([]domain.Order, error)
	GetOrderByID(id uint) (*domain.Order, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(r repository.OrderRepository) OrderService {
	return &orderService{repo: r}
}

func (s *orderService) CreateOrder(order domain.Order) (domain.Order, error) {
	return s.repo.Create(order)
}

func (s *orderService) UpdateOrder(order *domain.Order) error {
	return s.repo.Update(order)
}

func (s *orderService) DeleteOrder(id uint) error {
	return s.repo.Delete(id)
}

func (s *orderService) GetOrders() ([]domain.Order, error) {
	return s.repo.FindAll()
}

func (s *orderService) GetOrderByID(id uint) (*domain.Order, error) {
	return s.repo.FindByID(id)
}
