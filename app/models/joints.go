package models

import "time"

type JointsEvaluation struct {
	Id         uint `gorm:"primaryKey;autoIncrement"`
	PatientId  uint `gorm:"index"`
	RightAnkle int  `gorm:"not null"`
	LeftAnkle  int  `gorm:"not null"`
	RightKnee  int  `gorm:"not null"`
	LeftKnee   int  `gorm:"not null"`
	RightElbow int  `gorm:"not null"`
	LeftElbow  int  `gorm:"not null"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (JointsEvaluation) TableName() string {
	return "joints_evaluations"
}
