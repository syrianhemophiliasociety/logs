package models

import "time"

type Address struct {
	Id          uint   `gorm:"primaryKey;autoIncrement"`
	Governorate string `gorm:"index;not null"`
	Suburb      string `gorm:"index;not null"`
	Street      string `gorm:"index;not null"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (Address) TableName() string {
	return "addresses"
}
