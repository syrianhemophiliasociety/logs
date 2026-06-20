package models

import "time"

type VisitReason string

const (
	VisitReasonPrimaryProphylaxis   VisitReason = "primary_prophylaxis"
	VisitReasonSecondaryProphylaxis VisitReason = "secondary_prophylaxis"
	VisitReasonSurgery              VisitReason = "surgery"
	VisitReasonJointEvaluation      VisitReason = "joint_evaluation"
	VisitReasonJointInjection       VisitReason = "joint_injection"
	VisitReasonHemelibra            VisitReason = "hemelibra"
	VisitReasonTreatmentAtHome      VisitReason = "home_treatment"
	VisitReasonActiveBleeding       VisitReason = "active_bleeding"
)

type Visit struct {
	Id            uint        `gorm:"primaryKey;autoIncrement"`
	PatientId     uint        `gorm:"index;not null"`
	Reason        VisitReason `gorm:"not null"`
	Notes         string
	PatientWeight float64
	PatientHeight float64

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (Visit) TableName() string {
	return "visits"
}

type TreatmentDetails struct {
	Id          uint   `gorm:"primaryKey;autoIncrement"`
	Title       string `gorm:"not null"`
	ArabicTitle string `gorm:"not null"`
	Type        string `gorm:"not null"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (TreatmentDetails) TableName() string {
	return "treatment_details"
}

type PrescribedMedicine struct {
	Id                 uint `gorm:"primaryKey;autoIncrement"`
	VisitId            uint `gorm:"not null;index"`
	PatientId          uint `gorm:"not null"`
	Medicine           Medicine
	MedicineId         uint `gorm:"not null"`
	UsedAt             time.Time
	TreatmentDetails   TreatmentDetails `gorm:"-"`
	TreatmentDetailsId uint

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (PrescribedMedicine) TableName() string {
	return "prescribed_medicines"
}
