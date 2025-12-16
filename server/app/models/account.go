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

type AccountPermissions uint64

const (
	AccountPermissionReadAccounts AccountPermissions = 1 << iota
	AccountPermissionWriteAccounts
	AccountPermissionReadPatient
	AccountPermissionWritePatient
	AccountPermissionReadMedicine
	AccountPermissionWriteMedicine
	AccountPermissionReadVirus
	AccountPermissionWriteVirus
	AccountPermissionReadBloodTest
	AccountPermissionWriteBloodTest
	AccountPermissionReadOwnVisit
	AccountPermissionWriteOwnVisit
	AccountPermissionReadOtherVisits
	AccountPermissionWriteOtherVisits
	AccountPermissionReadDiagnoses
	AccountPermissionWriteDiagnoses
)

type Account struct {
	Id          uint               `gorm:"primaryKey;autoIncrement"`
	DisplayName string             `gorm:"not null"`
	Username    string             `gorm:"index;unique;not null"`
	Password    string             `gorm:"not null"`
	Type        AccountType        `gorm:"not null"`
	Permissions AccountPermissions `gorm:"not null"`

	CreatedAt time.Time `gorm:"index;not null"`
	UpdatedAt time.Time
}

func (Account) TableName() string {
	return "accounts"
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

func (a Account) HasPermission(p AccountPermissions) bool {
	return a.Permissions&p != 0
}
