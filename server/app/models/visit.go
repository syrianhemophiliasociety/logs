package models

import "time"

type VisitReason uint

const (
	VisitReasonSurgery VisitReason = 1 << iota
	VisitReasonJointEvaluation
	VisitReasonJointInjection
	VisitReasonHemelibra
	VisitReasonPafilaxes
)

type Visit struct {
	Id                  uint        `gorm:"primaryKey;autoIncrement"`
	PatientId           uint        `gorm:"not null"`
	Reason              VisitReason `gorm:"not null"`
	PrescribedMedicines []Medicine  `gorm:"many2many:prescribed_medicines;"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}
