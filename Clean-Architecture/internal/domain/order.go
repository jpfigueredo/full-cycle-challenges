package domain

import "time"

type Order struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	PatientID  int64     `json:"patient_id"`
	Medication string    `json:"medication"`
	Dosage     string    `json:"dosage"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
