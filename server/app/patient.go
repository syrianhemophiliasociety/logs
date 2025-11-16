package app

import (
	"shs/app/models"
	"time"
)

func (a *App) CreatePatient(patient models.Patient) (models.Patient, error) {
	return a.repo.CreatePatient(patient)
}

func (a *App) GetPatientById(id uint) (models.Patient, error) {
	patient, err := a.repo.GetPatientById(id)
	if err != nil {
		return models.Patient{}, err
	}

	viri, err := a.repo.ListViriForPatient(patient.Id)
	if err != nil {
		return models.Patient{}, err
	}

	patient.Viri = viri

	bloodTests, err := a.repo.ListPatientBloodTestResults(patient.Id)
	if err != nil {
		return models.Patient{}, err
	}

	patient.BloodTests = bloodTests

	return patient, nil
}

func (a *App) FindPatientsByVisitRange(from, to time.Time) ([]models.Patient, error) {
	return a.repo.FindPatientsByVisitDateRange(from, to)
}

func (a *App) FindPatientsByIndexFields(fields models.PatientIndexFields) ([]models.Patient, error) {
	return a.repo.FindPatientsByFields(fields)
}

func (a *App) ListMedicinesForVisit(visitId uint) ([]models.Medicine, error) {
	return a.repo.ListMedicinesForVisit(visitId)
}

func (a *App) UseMedicine(patientId, medicineId uint) error {
	return a.repo.UseMedicine(patientId, medicineId)
}
