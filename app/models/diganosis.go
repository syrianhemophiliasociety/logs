package models

import "time"

type Diagnosis struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	GroupName string `gorm:"not null"`
	Title     string `gorm:"not null"`
	// TODO: add filled fields

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (Diagnosis) TableName() string {
	return "diagnoses"
}

type DiagnosisResult struct {
	Id          uint `gorm:"primaryKey;autoIncrement"`
	DiagnosisId uint `gorm:"not null"`
	Diagnosis   Diagnosis
	PatientId   uint      `gorm:"index;not null"`
	DiagnosedAt time.Time `gorm:"not null"`
	// TODO: add filled fields

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (DiagnosisResult) TableName() string {
	return "diagnoses_results"
}
