package models

import (
	"time"

	"gorm.io/gorm"
)

type BlootTestUnit string

const (
	BlootTestUnitSecond     BlootTestUnit = "second"
	BlootTestUnitGram       BlootTestUnit = "g"
	BlootTestUnitPercentage BlootTestUnit = "%"
	BlootTestUnitCell       BlootTestUnit = "cell"
	BlootTestUnitML         BlootTestUnit = "mL"
	BlootTestUnitBU         BlootTestUnit = "BU"

	BlootTestUnitGramPerDeciLiter       BlootTestUnit = "g/dL"
	BlootTestUnitGramPerLiter           BlootTestUnit = "g/L"
	BlootTestUnitGramPerCubicCentimeter BlootTestUnit = "g/cm^3"

	BlootTestUnitCellPerCubicMilliLiter BlootTestUnit = "cell/mm^3"

	BlootTestUnitInternationalUnitPerDeciLiter BlootTestUnit = "IU/dL"
)

func BloodTestUnits() []BlootTestUnit {
	return []BlootTestUnit{
		BlootTestUnitSecond,
		BlootTestUnitGram,
		BlootTestUnitPercentage,
		BlootTestUnitCell,
		BlootTestUnitML,
		BlootTestUnitBU,
		BlootTestUnitGramPerDeciLiter,
		BlootTestUnitGramPerLiter,
		BlootTestUnitGramPerCubicCentimeter,
		BlootTestUnitCellPerCubicMilliLiter,
		BlootTestUnitInternationalUnitPerDeciLiter,
	}
}

type BloodTestField struct {
	Id          uint          `gorm:"primaryKey;autoIncrement"`
	BloodTestId uint          `gorm:"not null"`
	Name        string        `gorm:"not null"`
	Unit        BlootTestUnit `gorm:"not null"`
	MinValue    uint
	MaxValue    uint

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

type BloodTestFilledField struct {
	Id               uint `gorm:"primaryKey;autoIncrement"`
	BloodTestId      uint
	BloodTestFieldId uint
	ValueNumber      uint
	ValueString      string

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

type BloodTest struct {
	Id     uint             `gorm:"primaryKey;autoIncrement"`
	Name   string           `gorm:"not null"`
	Fields []BloodTestField `gorm:"foreignKey:BloodTestId"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
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

type BloodTestResult struct {
	Id           uint                   `gorm:"primaryKey;autoIncrement"`
	BloodTestId  uint                   `gorm:"not null"`
	PatientId    uint                   `gorm:"not null"`
	FilledFields []BloodTestFilledField `gorm:"foreignKey:BloodTestFieldId"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}
