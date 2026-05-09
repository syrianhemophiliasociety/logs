package models

import (
	"fmt"
	"time"
)

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

type PatientFirstVisitReason string

const (
	PatientFirstVisitReasonFamilyHistory PatientFirstVisitReason = "family_history"
	PatientFirstVisitReasonBleeding      PatientFirstVisitReason = "bleeding"
	PatientFirstVisitReasonReferral      PatientFirstVisitReason = "referral"
)

type Patient struct {
	Id                  uint                    `gorm:"primaryKey;autoIncrement"`
	PublicId            string                  `gorm:"index;not null;unique"`
	NationalId          string                  `gorm:"index;unique"`
	Nationality         string                  `gorm:"not null"`
	FirstName           string                  `gorm:"index;not null"`
	LastName            string                  `gorm:"index;not null"`
	FatherName          string                  `gorm:"index;not null"`
	MotherName          string                  `gorm:"index;not null"`
	PlaceOfBirth        Address                 `gorm:"not null"`
	PlaceOfBirthId      uint                    `gorm:"index;not null"`
	DateOfBirth         time.Time               `gorm:"not null"`
	Residency           Address                 `gorm:"not null"`
	ResidencyId         uint                    `gorm:"index;not null"`
	Gender              bool                    `gorm:"not null;index"`
	PhoneNumber         string                  `gorm:"index;not null"`
	FamilyHistoryExists bool                    `gorm:"not null"`
	FirstVisitReason    PatientFirstVisitReason `gorm:"not null"`
	BATScore            uint                    `gorm:"not null"`
	// TODO: keep only in the action's model
	Viruses           []Virus            `gorm:"many2many:has_viruses;"`
	BloodTestResults  []BloodTestResult  `gorm:"many2many:did_blood_tests;"`
	JointsEvaluations []JointsEvaluation `gorm:"many2many:patient_joint_evaluation;"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (Patient) TableName() string {
	return "patients"
}

func (p Patient) IndexId() string {
	return fmt.Sprintf("%s#%s#%s#%s", p.FirstName, p.LastName, p.FatherName, p.MotherName)
}

// FillEmptyFieldsUsingPublicId sets some empty fields with the value please_change_{publicId}
// it includes publicId since some fields are indexes.
// MUST ONLY BE USED WHEN USING THE IMPORT FUNCTIONALITY
func (p *Patient) FillEmptyFieldsUsingPublicId() {
	if p.NationalId == "" {
		p.NationalId = "please_change_" + p.PublicId
	}
	if p.PhoneNumber == "" {
		p.PhoneNumber = "please_change_" + p.PublicId
	}
	if p.PlaceOfBirth.Governorate == "" {
		p.PlaceOfBirth.Governorate = "please_change_" + p.PublicId
	}
	if p.PlaceOfBirth.Suburb == "" {
		p.PlaceOfBirth.Suburb = "please_change_" + p.PublicId
	}
	if p.Residency.Governorate == "" {
		p.Residency.Governorate = "please_change_" + p.PublicId
	}
	if p.Residency.Suburb == "" {
		p.Residency.Suburb = "please_change_" + p.PublicId
	}
}

type PatientUseMedicine struct {
	Id         uint `gorm:"primaryKey;autoIncrement"`
	PatientId  uint `gorm:"not null"`
	VisitId    uint `gorm:"not null"`
	MedicineId uint `gorm:"not null"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (PatientUseMedicine) TableName() string {
	return "patients_use_medicines"
}

type PatientId struct {
	Id       uint `gorm:"primaryKey;autoIncrement"`
	PublicId uint `gorm:"not null;index"`
}

func (PatientId) TableName() string {
	return "patient_ids"
}
