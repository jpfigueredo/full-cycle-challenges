package domain

import "time"

type Patient struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" binding:"required,min=3"`
	Age       int       `json:"age" binding:"gte=0,lte=130"`
	Email     string    `json:"email" binding:"required,email" gorm:"uniqueIndex"`
	BirthDate time.Time `json:"birth_date"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
