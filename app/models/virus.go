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

type HasVirus struct {
	Virus     Virus
	VirusId   uint `gorm:"primaryKey"`
	Patient   Patient
	PatientId uint `gorm:"primaryKey"`

	CreatedAt time.Time `gorm:"index;not null"`
}

func (HasVirus) TableName() string {
	return "has_viruses"
}
