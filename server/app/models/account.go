package models

import "time"

type AccountType string

const (
	AccountTypeSuperAdmin AccountType = "superadmin"
	AccountTypeAdmin      AccountType = "admin"
	AccountTypeSecritary  AccountType = "secritary"
	AccountTypePatient    AccountType = "patient"
)

type Account struct {
	Id          uint        `gorm:"primaryKey;autoIncrement"`
	DisplayName string      `gorm:"not null"`
	Username    string      `gorm:"index;unique;not null"`
	Password    string      `gorm:"not null"`
	Type        AccountType `gorm:"not null"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}
