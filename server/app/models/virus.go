package models

import "time"

type Virus struct {
	Id   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"not null"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}
