package models

import (
	"time"

	"gorm.io/gorm"
)

type BlootTestUnit string

const (
	BlootTestUnitSecond BlootTestUnit = "second"
	BlootTestUnitMinute BlootTestUnit = "minute"

	BlootTestUnitPercentage BlootTestUnit = "%"
	BlootTestUnitCell       BlootTestUnit = "cell"
	BlootTestUnitBU         BlootTestUnit = "BU"

	BlootTestUnitGram                   BlootTestUnit = "g"
	BlootTestUnitPicoGram               BlootTestUnit = "pg"
	BlootTestUnitGramPerDeciLiter       BlootTestUnit = "g/dL"
	BlootTestUnitGramPerLiter           BlootTestUnit = "g/L"
	BlootTestUnitGramPerCubicCentimeter BlootTestUnit = "g/cm^3"
	BlootTestUnitMilligramPerDeciLiter  BlootTestUnit = "mg/dL"
	BlootTestUnitMicroUnitPerMilliLiter BlootTestUnit = "mcU/mL"
	BlootTestUnitMicroGramPerDeciLiter  BlootTestUnit = "mcg/dL"
	BlootTestUnitNanoGramPerDeciLiter   BlootTestUnit = "ng/dL"
	BlootTestUnitPicoGramPerDeciLiter   BlootTestUnit = "pg/dL"

	BlootTestUnitML         BlootTestUnit = "mL"
	BlootTestUnitFemtoLiter BlootTestUnit = "fL"

	BlootTestUnitInternationalUnitPerDeciLiter BlootTestUnit = "IU/dL"
	BlootTestUnitUnitPerLiter                  BlootTestUnit = "U/L"

	BlootTestUnitCellPerCubicMilliLiter         BlootTestUnit = "cell/mm^3"
	BlootTestUnitThousandCellPerCubicMillimeter BlootTestUnit = "10^3 cell/mm^3"
	BlootTestUnitMillionCellPerCubicMillimeter  BlootTestUnit = "10^6 cell/mm^3"

	BlootTestUnitRatioOrIndex BlootTestUnit = "-"
)

func BloodTestUnits() []BlootTestUnit {
	return []BlootTestUnit{
		BlootTestUnitSecond,
		BlootTestUnitMinute,
		BlootTestUnitPercentage,
		BlootTestUnitCell,
		BlootTestUnitBU,
		BlootTestUnitGram,
		BlootTestUnitPicoGram,
		BlootTestUnitGramPerDeciLiter,
		BlootTestUnitGramPerLiter,
		BlootTestUnitGramPerCubicCentimeter,
		BlootTestUnitMilligramPerDeciLiter,
		BlootTestUnitMicroUnitPerMilliLiter,
		BlootTestUnitMicroGramPerDeciLiter,
		BlootTestUnitNanoGramPerDeciLiter,
		BlootTestUnitPicoGramPerDeciLiter,
		BlootTestUnitML,
		BlootTestUnitFemtoLiter,
		BlootTestUnitInternationalUnitPerDeciLiter,
		BlootTestUnitUnitPerLiter,
		BlootTestUnitCellPerCubicMilliLiter,
		BlootTestUnitThousandCellPerCubicMillimeter,
		BlootTestUnitMillionCellPerCubicMillimeter,
		BlootTestUnitRatioOrIndex,
	}
}

type BloodTestField struct {
	Id             uint          `gorm:"primaryKey;autoIncrement"`
	BloodTestId    uint          `gorm:"not null"`
	Name           string        `gorm:"not null"`
	Unit           BlootTestUnit `gorm:"not null"`
	MinValueNumber float64
	MinValueString string
	MaxValueNumber float64
	MaxValueString string

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (BloodTestField) TableName() string {
	return "blood_test_fields"
}

type BloodTest struct {
	Id     uint             `gorm:"primaryKey;autoIncrement"`
	Name   string           `gorm:"not null"`
	Fields []BloodTestField `gorm:"foreignKey:BloodTestId"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (BloodTest) TableName() string {
	return "blood_tests"
}

func (bt *BloodTest) AfterDelete(tx *gorm.DB) error {
	for i := range bt.Fields {
		err := tx.
			Model(new(BloodTestField)).
			Delete(&bt.Fields[i], "id = ?", bt.Fields[i].Id).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}

type BloodTestFilledField struct {
	Id                uint `gorm:"primaryKey;autoIncrement"`
	BloodTestResultId uint
	BloodTestFieldId  uint
	ValueNumber       float64
	ValueString       string

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (BloodTestFilledField) TableName() string {
	return "blood_test_filled_fields"
}

type BloodTestResult struct {
	Id           uint `gorm:"primaryKey;autoIncrement"`
	BloodTestId  uint `gorm:"not null"`
	BloodTest    BloodTest
	PatientId    uint                   `gorm:"index;not null"`
	FilledFields []BloodTestFilledField `gorm:"foreignKey:BloodTestResultId"`
	Pending      bool                   `gorm:"not null"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (BloodTestResult) TableName() string {
	return "blood_test_results"
}
