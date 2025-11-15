package models

import "time"

type Medicine struct {
	Id   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"not null"`
	Dose int    `gorm:"not null"`
	Unit string `gorm:"not null"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}
