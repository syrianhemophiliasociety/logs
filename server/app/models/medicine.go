package models

import "time"

type MedicineDoseUnit string

const (
	MedicineDoseUnitMilliLiter MedicineDoseUnit = "mL"
	MedicineDoseUnitLiter      MedicineDoseUnit = "L"
	MedicineDoseUnitGram       MedicineDoseUnit = "g"
)

type Medicine struct {
	Id   uint             `gorm:"primaryKey;autoIncrement"`
	Name string           `gorm:"not null"`
	Dose int              `gorm:"not null"`
	Unit MedicineDoseUnit `gorm:"not null"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}
