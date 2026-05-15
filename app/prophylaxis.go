package app

import (
	"shs/app/models"
	"time"
)

func (a *App) CreateProphylaxis(pp models.Prophylaxis) (models.Prophylaxis, error) {
	return a.repo.CreateProphylaxis(pp)
}

func (a *App) ListProphylaxesForPatient(patientId uint) ([]models.Prophylaxis, error) {
	return a.repo.ListProphylaxesForPatient(patientId)
}

func (a *App) DeleteProphylaxisForPatient(id, patientId uint) error {
	return a.repo.DeleteProphylaxisForPatient(id, patientId)
}

func (a *App) SetProphylaxisEndDateForPatient(id, patientId uint, endDate time.Time) (models.Prophylaxis, error) {
	return a.repo.SetProphylaxisEndDateForPatient(id, patientId, endDate)
}

func (a *App) SetProphylaxisChosenForPatient(id, patientId uint, chosen bool) (models.Prophylaxis, error) {
	return a.repo.SetProphylaxisChosenForPatient(id, patientId, chosen)
}
