package repository

import (
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/domain"
	"gorm.io/gorm"
)

type PatientRepository interface {
	FindAll() ([]domain.Patient, error)
	FindByID(id uint) (*domain.Patient, error)
	Create(patient domain.Patient) (domain.Patient, error)
	Update(patient *domain.Patient) error
	Delete(id uint) error
}

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) FindAll() ([]domain.Patient, error) {
	var patients []domain.Patient
	if err := r.db.Find(&patients).Error; err != nil {
		return nil, err
	}
	return patients, nil
}

func (r *patientRepository) FindByID(id uint) (*domain.Patient, error) {
	var patient domain.Patient
	if err := r.db.First(&patient, id).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *patientRepository) Create(patient domain.Patient) (domain.Patient, error) {
	if err := r.db.Create(&patient).Error; err != nil {
		return domain.Patient{}, err
	}
	return patient, nil
}

func (r *patientRepository) Update(patient *domain.Patient) error {
	return r.db.Save(patient).Error
}

func (r *patientRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Patient{}, id).Error
}
