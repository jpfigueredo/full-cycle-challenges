package service

import (
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/domain"
	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/repository"
)

type PatientService interface {
	GetPatients() ([]domain.Patient, error)
	GetPatientByID(id uint) (*domain.Patient, error)
	CreatePatient(patient domain.Patient) (domain.Patient, error)
	UpdatePatient(patient *domain.Patient) error
	DeletePatient(id uint) error
}

type patientService struct {
	repo repository.PatientRepository
}

func NewPatientService(r repository.PatientRepository) PatientService {
	return &patientService{repo: r}
}

func (s *patientService) GetPatients() ([]domain.Patient, error) {
	return s.repo.FindAll()
}

func (s *patientService) GetPatientByID(id uint) (*domain.Patient, error) {
	return s.repo.FindByID(id)
}

func (s *patientService) CreatePatient(patient domain.Patient) (domain.Patient, error) {
	return s.repo.Create(patient)
}

func (s *patientService) UpdatePatient(patient *domain.Patient) error {
	return s.repo.Update(patient)
}

func (s *patientService) DeletePatient(id uint) error {
	return s.repo.Delete(id)
}
