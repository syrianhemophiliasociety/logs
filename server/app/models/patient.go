package models

import "time"

type PatientIndexFields struct {
	PublicId     string
	NationalId   string
	FirstName    string
	LastName     string
	FatherName   string
	MotherName   string
	PlaceOfBirth Address
	Residency    Address
	PhoneNumber  string
}

type Patient struct {
	Id             uint              `gorm:"primaryKey;autoIncrement"`
	PublicId       string            `gorm:"index;not null;unique"`
	NationalId     string            `gorm:"index;not null;unique"`
	Nationality    string            `gorm:"not null"`
	FirstName      string            `gorm:"index;not null;unique"`
	LastName       string            `gorm:"index;not null;unique"`
	FatherName     string            `gorm:"index;not null;unique"`
	MotherName     string            `gorm:"index;not null;unique"`
	PlaceOfBirth   Address           `gorm:"not null"`
	PlaceOfBirthId uint              `gorm:"index;not null"`
	DateOfBirth    time.Time         `gorm:"not null"`
	Residency      Address           `gorm:"not null"`
	ResidencyId    uint              `gorm:"index;not null"`
	Gender         bool              `gorm:"not null"`
	PhoneNumber    string            `gorm:"index;not null"`
	BATScore       uint              `gorm:"not null"`
	Viri           []Virus           `gorm:"not null;many2many:has_viruses;"`
	BloodTests     []BloodTestResult `gorm:"not null;many2many:did_blood_tests;"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

type PatientUseMedicine struct {
	Id         uint `gorm:"primaryKey;autoIncrement"`
	PatientId  uint `gorm:"not null"`
	VisitId    uint `gorm:"not null"`
	MedicineId uint `gorm:"not null"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}
