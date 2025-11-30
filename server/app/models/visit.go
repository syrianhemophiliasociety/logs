package models

import "time"

type VisitReason string

const (
	VisitReasonSurgery         VisitReason = "surgery"
	VisitReasonJointEvaluation VisitReason = "joint_evaluation"
	VisitReasonJointInjection  VisitReason = "joint_injection"
	VisitReasonHemelibra       VisitReason = "hemelibra"
	VisitReasonPafilaxes       VisitReason = "pafilaxes"
)

type Visit struct {
	Id        uint        `gorm:"primaryKey;autoIncrement"`
	PatientId uint        `gorm:"not null"`
	Reason    VisitReason `gorm:"not null"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
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
