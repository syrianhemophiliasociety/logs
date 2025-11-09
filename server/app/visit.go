package app

import "shs/app/models"

func (a *App) CreatePatientVisit(visit models.Visit) (models.Visit, error) {
	return a.repo.CreatePatientVisit(visit)
}
