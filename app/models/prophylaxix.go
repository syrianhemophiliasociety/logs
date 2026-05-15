package models

import "time"

type Prophylaxis struct {
	Id               uint `gorm:"primaryKey;autoIncrement"`
	PatientId        uint `gorm:"index;not null"`
	MedicineId       uint `gorm:"index;not null"`
	Medicine         Medicine
	MedicineDose     int     `gorm:"not null"`
	Title            string  `gorm:"not null"`
	FrequencyPerDays float32 `gorm:"not null"`
	EndDate          time.Time
	Chosen           bool

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (Prophylaxis) TableName() string {
	return "prophylaxes"
}
