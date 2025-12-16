package app

import "shs/app/models"

func (a *App) CreateVirus(virus models.Virus) (models.Virus, error) {
	return a.repo.CreateVirus(virus)
}

func (a *App) DeleteVirus(id uint) error {
	return a.repo.DeleteVirus(id)
}

func (a *App) ListAllViruses() ([]models.Virus, error) {
	return a.repo.ListAllViruses()
}
