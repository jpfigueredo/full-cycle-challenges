package repository

import (
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/domain"
	"gorm.io/gorm"
)

type OrderRepository interface {
	FindAll() ([]domain.Order, error)
	FindByID(id uint) (*domain.Order, error)
	Create(order *domain.Order) error
	Update(order *domain.Order) error
	Delete(id uint) error
}

type GormOrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *GormOrderRepository {
	return &GormOrderRepository{db: db}
}

func (r *GormOrderRepository) FindAll() ([]domain.Order, error) {
	var orders []domain.Order
	if err := r.db.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *GormOrderRepository) FindByID(id uint) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *GormOrderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *GormOrderRepository) Update(order *domain.Order) error {
	return r.db.Save(order).Error
}

func (r *GormOrderRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Order{}, id).Error
}
