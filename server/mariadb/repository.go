package mariadb

import (
	"errors"
	"fmt"
	"shs/app"
	"shs/app/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	client *gorm.DB
}

func New() (*Repository, error) {
	conn, err := dbConnector()
	if err != nil {
		return nil, err
	}

	return &Repository{
		client: conn,
	}, nil
}

// --------------------------------
// App Repository
// --------------------------------

func (r *Repository) GetAccount(id uint) (models.Account, error) {
	var account models.Account

	err := tryWrapDbError(
		r.client.
			Model(new(models.Account)).
			First(&account, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.Account{}, &app.ErrNotFound{
			ResourceName: "account",
		}
	}
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (r *Repository) GetAccountByUsername(username string) (models.Account, error) {
	var account models.Account

	err := tryWrapDbError(
		r.client.
			Model(new(models.Account)).
			First(&account, "username = ?", username).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.Account{}, &app.ErrNotFound{
			ResourceName: "account",
		}
	}
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (r *Repository) CreateAccount(account models.Account) (models.Account, error) {
	account.CreatedAt = time.Now().UTC()
	account.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.Account)).
			Create(&account).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.Account{}, &app.ErrExists{
			ResourceName: "account",
		}
	}
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (r *Repository) CreateBloodTest(bt models.BloodTest) (models.BloodTest, error) {
	bt.CreatedAt = time.Now().UTC()
	bt.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.BloodTest)).
			Create(&bt).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.BloodTest{}, &app.ErrExists{
			ResourceName: "blood_test",
		}
	}
	if err != nil {
		return models.BloodTest{}, err
	}

	return bt, nil
}

func (r *Repository) DeleteBloodTest(id uint) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.BloodTest)).
			Delete(&models.BloodTest{Id: id}, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "blood_test",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetBloodTest(id uint) (models.BloodTest, error) {
	var bt models.BloodTest

	err := tryWrapDbError(
		r.client.
			Model(new(models.BloodTest)).
			Preload("Fields").
			First(&bt, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.BloodTest{}, &app.ErrNotFound{
			ResourceName: "blood_test",
		}
	}
	if err != nil {
		return models.BloodTest{}, err
	}

	return bt, nil
}

func (r *Repository) UpdateBloodTest(id uint, bt models.BloodTest) (models.BloodTest, error) {
	return models.BloodTest{}, errors.New("not implemented")
}

func (r *Repository) ListAllBloodTests() ([]models.BloodTest, error) {
	var bloodTests []models.BloodTest

	err := tryWrapDbError(
		r.client.
			Model(new(models.BloodTest)).
			Preload("Fields").
			Find(&bloodTests).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return nil, &app.ErrNotFound{
			ResourceName: "blood_test",
		}
	}
	if err != nil {
		return nil, err
	}

	return bloodTests, nil
}

func (r *Repository) CreateVirus(virus models.Virus) (models.Virus, error) {
	virus.CreatedAt = time.Now().UTC()
	virus.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.Virus)).
			Create(&virus).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.Virus{}, &app.ErrExists{
			ResourceName: "virus",
		}
	}
	if err != nil {
		return models.Virus{}, err
	}

	return virus, nil

}

func (r *Repository) DeleteVirus(id uint) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Virus)).
			Delete(&models.Virus{Id: id}, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "virus",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListAllViri() ([]models.Virus, error) {
	var viri []models.Virus

	err := tryWrapDbError(
		r.client.
			Model(new(models.Virus)).
			Find(&viri).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return nil, &app.ErrNotFound{
			ResourceName: "virus",
		}
	}
	if err != nil {
		return nil, err
	}

	return viri, nil
}

func (r *Repository) CreateMedicine(medicine models.Medicine) (models.Medicine, error) {
	medicine.CreatedAt = time.Now().UTC()
	medicine.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.Medicine)).
			Create(&medicine).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.Medicine{}, &app.ErrExists{
			ResourceName: "medicine",
		}
	}
	if err != nil {
		return models.Medicine{}, err
	}

	return medicine, nil

}

func (r *Repository) DeleteMedicine(id uint) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Medicine)).
			Delete(&models.Medicine{Id: id}, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "medicine",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListAllMedicines() ([]models.Medicine, error) {
	var viri []models.Medicine

	err := tryWrapDbError(
		r.client.
			Model(new(models.Medicine)).
			Find(&viri).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return nil, &app.ErrNotFound{
			ResourceName: "medicine",
		}
	}
	if err != nil {
		return nil, err
	}

	return viri, nil
}

func (r *Repository) CreatePatient(patient models.Patient) (models.Patient, error) {
	patient.CreatedAt = time.Now().UTC()
	patient.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.Patient)).
			Create(&patient).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.Patient{}, &app.ErrExists{
			ResourceName: "patient",
		}
	}
	if err != nil {
		return models.Patient{}, err
	}

	return patient, nil
}

func (r *Repository) FindPatientsByVisitDateRange(from, to time.Time) ([]models.Patient, error) {
	return nil, errors.New("not inmplemented")
}

func (r *Repository) FindPatientsByFields(patientIndexFields models.PatientIndexFields) ([]models.Patient, error) {
	return nil, errors.New("not implemented")
}

func (r *Repository) ListMedicinesForVisit(visitId uint) ([]models.Medicine, error) {
	return nil, errors.New("not implemented")
}

func (r *Repository) UseMedicine(patientId, medicineId uint) error {
	return errors.New("not implemented")
}

func (r *Repository) CreatePatientVisit(visit models.Visit) (models.Visit, error) {
	visit.CreatedAt = time.Now().UTC()
	visit.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.Visit)).
			Create(&visit).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.Visit{}, &app.ErrExists{
			ResourceName: "visit",
		}
	}
	if err != nil {
		return models.Visit{}, err
	}

	return visit, nil
}

func (r *Repository) CreateAddress(address models.Address) (models.Address, error) {
	address.CreatedAt = time.Now().UTC()
	address.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.Address)).
			Create(&address).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.Address{}, &app.ErrExists{
			ResourceName: "address",
		}
	}
	if err != nil {
		return models.Address{}, err
	}

	return address, nil
}

func (r *Repository) GetAllAddresses() ([]models.Address, error) {
	var addresses []models.Address

	err := tryWrapDbError(
		r.client.
			Model(new(models.Address)).
			Find(&addresses).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return nil, &app.ErrNotFound{
			ResourceName: "address",
		}
	}
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *Repository) GetAllAddressesALike(searchAddress models.Address) ([]models.Address, error) {
	findQuery := make([]string, 0, 3)
	findArgs := make([]any, 0, 3)
	if searchAddress.Governorate != "" {
		findQuery = append(findQuery, "LOWER(governorate) LIKE LOWER(?)")
		findArgs = append(findArgs, likeArg(searchAddress.Governorate))
	}
	if searchAddress.Suburb != "" {
		findQuery = append(findQuery, "LOWER(suburb) LIKE LOWER(?)")
		findArgs = append(findArgs, likeArg(searchAddress.Suburb))
	}
	if searchAddress.Street != "" {
		findQuery = append(findQuery, "LOWER(street) LIKE LOWER(?)")
		findArgs = append(findArgs, likeArg(searchAddress.Street))
	}

	var addresses []models.Address

	err := tryWrapDbError(
		r.client.
			Model(new(models.Address)).
			Where(strings.Join(findQuery, " AND "), findArgs...).
			Find(&addresses).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return nil, &app.ErrNotFound{
			ResourceName: "address",
		}
	}
	if err != nil {
		return nil, err
	}

	return addresses, nil

}

func (r *Repository) DeleteAddress(id uint) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Address)).
			Delete(&models.Address{Id: id}, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "address",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func likeArg(arg string) string {
	return fmt.Sprintf("%%%s%%", arg)
}
