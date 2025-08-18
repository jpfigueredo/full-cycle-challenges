package domain

import "time"

type Order struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Item       string    `json:"item" binding:"required,min=3"`
	Amount     int       `json:"amount" binding:"required,gt=0"`
	PatientID  int64     `json:"patient_id"`
	Medication string    `json:"medication"`
	Dosage     string    `json:"dosage"`
	Status     string    `json:"status" binding:"required,oneof=OPEN CLOSED CANCELLED"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
