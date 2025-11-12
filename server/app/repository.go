package app

import (
	"shs/app/models"
	"time"
)

type Repository interface {
	GetAccount(id uint) (models.Account, error)
	GetAccountByUsername(username string) (models.Account, error)
	CreateAccount(account models.Account) (models.Account, error)
	ListAllAccounts(types []models.AccountType) ([]models.Account, error)
	DeleteAccount(id uint) error

	CreateBloodTest(bt models.BloodTest) (models.BloodTest, error)
	DeleteBloodTest(id uint) error
	GetBloodTest(id uint) (models.BloodTest, error)
	UpdateBloodTest(id uint, bt models.BloodTest) (models.BloodTest, error)
	ListAllBloodTests() ([]models.BloodTest, error)

	CreateVirus(virus models.Virus) (models.Virus, error)
	DeleteVirus(id uint) error
	ListAllViri() ([]models.Virus, error)

	CreateMedicine(medicine models.Medicine) (models.Medicine, error)
	DeleteMedicine(id uint) error
	ListAllMedicines() ([]models.Medicine, error)

	CreatePatient(patient models.Patient) (models.Patient, error)
	FindPatientsByVisitDateRange(from, to time.Time) ([]models.Patient, error)
	FindPatientsByFields(patientIndexFields models.PatientIndexFields) ([]models.Patient, error)

	ListMedicinesForVisit(visitId uint) ([]models.Medicine, error)
	UseMedicine(patientId, medicineId uint) error

	CreatePatientVisit(visit models.Visit) (models.Visit, error)

	CreateAddress(address models.Address) (models.Address, error)
	GetAllAddresses() ([]models.Address, error)
	GetAllAddressesALike(searchAddress models.Address) ([]models.Address, error)
	DeleteAddress(id uint) error
}
