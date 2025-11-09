package app

import (
	"shs/app/models"
	"time"
)

func (a *App) CreatePatient(patient models.Patient) (models.Patient, error) {
	return a.repo.CreatePatient(patient)
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
