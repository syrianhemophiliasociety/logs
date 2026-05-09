package models

import "time"

type Virus struct {
	Id                    uint        `gorm:"primaryKey;autoIncrement"`
	Name                  string      `gorm:"not null"`
	IdentifyingBloodTests []BloodTest `gorm:"not null;many2many:identifying_blood_tests;"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (Virus) TableName() string {
	return "viruses"
}
