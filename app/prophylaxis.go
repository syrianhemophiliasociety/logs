package app

import "shs/app/models"

func (a *App) CreateProphylaxis(pp models.Prophylaxis) (models.Prophylaxis, error) {
	return a.repo.CreateProphylaxis(pp)
}

func (a *App) ListProphylaxesForPatient(patientId uint) ([]models.Prophylaxis, error) {
	return a.repo.ListProphylaxesForPatient(patientId)
}
