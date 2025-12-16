package models

import "time"

type Medicine struct {
	Id           uint      `gorm:"primaryKey;autoIncrement"`
	Name         string    `gorm:"not null"`
	Dose         int       `gorm:"not null"`
	Unit         string    `gorm:"not null"`
	Amount       int       `gorm:"not null"`
	ExpiresAt    time.Time `gorm:"not null"`
	ReceivedAt   time.Time `gorm:"not null"`
	Manufacturer string    `gorm:"not null"`
	BatchNumber  string    `gorm:"not null"`
	FactorType   string    `gorm:"not null"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (Medicine) TableName() string {
	return "medicines"
}
