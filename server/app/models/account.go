package models

import (
	"errors"
	"slices"
	"time"
)

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

func (a Account) CheckType(accountTypes ...AccountType) error {
	if a.Type == AccountTypeSuperAdmin {
		return nil
	}

	if slices.Contains(accountTypes, a.Type) {
		return nil
	}

	return errors.New("invalid account type")
}
