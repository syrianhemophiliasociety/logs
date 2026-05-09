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

type PrescribedMedicine struct {
	Id         uint `gorm:"primaryKey;autoIncrement"`
	VisitId    uint `gorm:"not null;index"`
	PatientId  uint `gorm:"not null"`
	MedicineId uint `gorm:"not null"`
	UsedAt     time.Time

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (PrescribedMedicine) TableName() string {
	return "prescribed_medicines"
}
