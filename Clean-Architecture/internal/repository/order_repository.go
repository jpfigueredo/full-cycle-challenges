package repository

import (
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/domain"
	"gorm.io/gorm"
)

type OrderRepository interface {
	FindAll() ([]domain.Order, error)
	FindByID(id uint) (*domain.Order, error)
	Create(order domain.Order) (domain.Order, error)
	Update(order *domain.Order) error
	Delete(id uint) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) FindAll() ([]domain.Order, error) {
	var orders []domain.Order
	if err := r.db.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) FindByID(id uint) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) Create(order domain.Order) (domain.Order, error) {
	if err := r.db.Create(&order).Error; err != nil {
		return domain.Order{}, err
	}
	return order, nil
}

func (r *orderRepository) Update(order *domain.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Order{}, id).Error
}
