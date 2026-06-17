package app

import (
	"shs/app/models"
	"time"
)

func (a *App) CreatePatientVisit(visit models.Visit) (models.Visit, error) {
	return a.repo.CreatePatientVisit(visit)
}

func (a *App) CreatePrescribedMedicine(pm models.PrescribedMedicine) (models.PrescribedMedicine, error) {
	return a.repo.CreatePrescribedMedicine(pm)
}

func (a *App) ListVisitsOnTimeRange(from, to time.Time) ([]models.Visit, error) {
	return a.repo.ListVisitsOnTimeRange(from, to)
}

func (a *App) ListAllTreatmentDetails() ([]models.TreatmentDetails, error) {
	return a.repo.ListAllTreatmentDetails()
}
